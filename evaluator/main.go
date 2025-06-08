package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strings"

	"github.com/akhmadst1/tugas-akhir-backend/internal/services"
	"github.com/akhmadst1/tugas-akhir-backend/pkg"
	"github.com/joho/godotenv"
)

type ChatDisharmony struct {
	Result   bool   `json:"result"`
	Analysis string `json:"analysis"`
}

// CosineSimilarity calculates similarity between two float slices
func CosineSimilarity(vec1, vec2 []float64) float64 {
	if len(vec1) != len(vec2) {
		return 0.0
	}
	var dotProduct, normA, normB float64
	for i := range vec1 {
		dotProduct += vec1[i] * vec2[i]
		normA += vec1[i] * vec1[i]
		normB += vec2[i] * vec2[i]
	}
	if normA == 0 || normB == 0 {
		return 0.0
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// Tokenize splits and normalizes a string (lowercase, no punctuation)
func Tokenize(text string) []string {
	text = strings.ToLower(text)
	re := regexp.MustCompile(`[^\w\s]`)
	text = re.ReplaceAllString(text, "")
	return strings.Fields(text)
}

// ComputeF1 calculates precision, recall, and F1 based on token overlap
func ComputeF1Score(prediction, reference string) (precision, recall, f1 float64) {
	predTokens := Tokenize(prediction)
	refTokens := Tokenize(reference)

	if len(predTokens) == 0 || len(refTokens) == 0 {
		return 0, 0, 0
	}

	refCounts := make(map[string]int)
	for _, t := range refTokens {
		refCounts[t]++
	}

	overlap := 0
	for _, t := range predTokens {
		if refCounts[t] > 0 {
			overlap++
			refCounts[t]--
		}
	}

	precision = float64(overlap) / float64(len(predTokens))
	recall = float64(overlap) / float64(len(refTokens))
	if precision+recall > 0 {
		f1 = 2 * precision * recall / (precision + recall)
	}

	return
}

// EvaluateDisharmonySimilarity compares ground truth and LLM response
func EvaluateDisharmonySimilarity(groundTruth string, llmOutput string) (float64, error) {
	gtEmb, err := pkg.GetGeminiEmbedding(groundTruth)
	if err != nil {
		return 0.0, err
	}

	llmEmb, err := pkg.GetGeminiEmbedding(llmOutput)
	if err != nil {
		return 0.0, err
	}

	similarity := CosineSimilarity(gtEmb, llmEmb)
	return similarity, nil
}

func ExtractChatDisharmony(llmOutput string) (ChatDisharmony, error) {
	// Remove any code block formatting like ```json or ``` from LLM output
	cleaned := strings.TrimSpace(llmOutput)
	if strings.HasPrefix(cleaned, "```json") {
		cleaned = strings.TrimPrefix(cleaned, "```json")
	} else if strings.HasPrefix(cleaned, "```") {
		cleaned = strings.TrimPrefix(cleaned, "```")
	}
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	var parsed ChatDisharmony
	if err := json.Unmarshal([]byte(cleaned), &parsed); err != nil {
		return ChatDisharmony{}, errors.New("failed to parse JSON from LLM output")
	}
	return parsed, nil
}

// BatchEvaluate compares LLM responses for all test cases
func BatchEvaluate(testCaseFile string, filename string) error {
	testCases, err := pkg.LoadTestCases(testCaseFile)
	if err != nil {
		return err
	}

	var output bytes.Buffer

	var totalSimilarity, totalPrecision, totalRecall, totalF1 float64
	var successfulCases, correctClassification, incorrectClassification int

	// Generate LLM output for each test case
	for _, testCase := range testCases {
		var fullRegulationText string
		for _, reg := range testCase.Regulations {
			fullRegulationText += fmt.Sprintf("Document: %s\nArticle: %s\nContent: %s\n\n", reg.Document, reg.Article, reg.Content)
		}
		var promptBuilder strings.Builder

		// ** -------------------- ZERO SHOT --------------------
		// promptBuilder.WriteString(pkg.ZeroShot(fullRegulationText))

		// ** -------------------- FEW SHOT --------------------
		// promptBuilder.WriteString(pkg.FewShot(fullRegulationText, testCase.ID))

		// ** ----------------------- CHAIN OF THOUGHT -------------------
		promptBuilder.WriteString(pkg.ChainOfThought(fullRegulationText))

		// ** ----------------------- FEW SHOT + CHAIN OF THOUGHT -------------------
		// promptBuilder.WriteString(pkg.FewShotChainOfThought(fullRegulationText, testCase.ID))

		// ------------------------ END OF PROMPT ------------------------

		// OPENAI API call
		llmOutput, err := services.EvaluateOpenAIDisharmonyAnalysis(promptBuilder.String())

		// GEMINI API call
		// llmOutput, err := services.EvaluateGeminiDisharmonyAnalysis(promptBuilder.String())

		// LLAMA API call
		// llmOutput, err := services.EvaluateLlamaDisharmonyAnalysis(promptBuilder.String())

		// HUGGING FACE LLAMA API call
		// llmOutput, err := services.EvaluateHuggingFaceLLaMADisharmonyAnalysis(promptBuilder.String())

		if err != nil {
			msg := fmt.Sprintf("Error generating LLM output for test case %s: %v\n", testCase.ID, err)
			fmt.Print(msg)
			output.WriteString(msg)
			continue
		}

		parsed, err := ExtractChatDisharmony(llmOutput)
		if err != nil {
			msg := fmt.Sprintf("Error parsing JSON response for test case %s: %v\n", testCase.ID, err)
			fmt.Print(msg)
			output.WriteString(msg)
			continue
		}

		// Compare LLM output with ground truth disharmony
		similarity, err := EvaluateDisharmonySimilarity(testCase.Disharmony, parsed.Analysis)
		if err != nil {
			msg := fmt.Sprintf("Error comparing test case %s: %v\n", testCase.ID, err)
			fmt.Print(msg)
			output.WriteString(msg)
			continue
		}

		p, r, f1 := ComputeF1Score(testCase.Disharmony, parsed.Analysis)

		totalSimilarity += similarity
		totalPrecision += p
		totalRecall += r
		totalF1 += f1
		successfulCases++
		if parsed.Result {
			correctClassification++
		} else {
			incorrectClassification++
		}

		// Build result
		result := fmt.Sprintf(
			"--------------------------------------------------\n"+
				"Test Case: %s - %s\n"+
				"Ground Truth Disharmony:\n%s\n\n"+
				"Disharmony Result: %t\n"+
				"Disharmony Analysis:\n%s\n\n"+
				"Similarity Score: %.4f\n"+
				"Precision: %.4f\n"+
				"Recall: %.4f\n"+
				"F1 Score: %.4f\n\n",
			testCase.ID, testCase.Title, testCase.Disharmony, parsed.Result, parsed.Analysis, similarity, p, r, f1,
		)
		result += ("--------------------------------------------------\n\n")

		output.WriteString(result)
	}

	if successfulCases > 0 {
		avgSimilarity := totalSimilarity / float64(successfulCases)
		avgPrecision := totalPrecision / float64(successfulCases)
		avgRecall := totalRecall / float64(successfulCases)
		avgF1 := totalF1 / float64(successfulCases)

		summary := fmt.Sprintf(
			"\n==================== Overall Metrics Summary ====================\n"+
				"Total Cases Evaluated: %d\n"+
				"True: %d\n"+
				"False: %d\n"+
				"Average Similarity Score: %.4f\n"+
				"Average Precision:        %.4f\n"+
				"Average Recall:           %.4f\n"+
				"Average F1 Score:         %.4f\n"+
				"=================================================================\n",
			successfulCases, correctClassification, incorrectClassification, avgSimilarity, avgPrecision, avgRecall, avgF1,
		)
		output.WriteString(summary)
	}

	// Save output to file
	err = os.WriteFile(filename, output.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing results to file: %v", err)
	}

	fmt.Printf("Results saved to %s", filename)
	return nil
}

func main() {
	// Load .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Default output filename
	outputFile := "../result/default-name-result.txt"

	// Check for command-line argument
	if len(os.Args) > 1 {
		outputFile = "../result/" + os.Args[1] + ".txt"
	}

	// Run the batch evaluation
	err = BatchEvaluate("../data/test_cases_14.json", outputFile)
	if err != nil {
		fmt.Println("Error during evaluation:", err)
	}
}

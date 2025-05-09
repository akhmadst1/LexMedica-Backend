package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strings"

	"github.com/akhmadst1/tugas-akhir-backend/internal/services"
	"github.com/akhmadst1/tugas-akhir-backend/pkg"
	"github.com/joho/godotenv"
)

// EmbeddingRequest is the payload for OpenAI embeddings API
type EmbeddingRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

// EmbeddingResponse holds the API response
type EmbeddingResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
}

// Regulation represents each legal regulation
type Regulation struct {
	Document string `json:"document"`
	Article  string `json:"article"`
	Content  string `json:"content"`
}

// TestCase represents the structure of each test case
type TestCase struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Regulations []Regulation `json:"regulations"`
	Disharmony  string       `json:"disharmony"`
}

// GetEmbedding fetches embedding vector from OpenAI API
func GetEmbedding(text string) ([]float64, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	url := "https://api.openai.com/v1/embeddings"

	payload := EmbeddingRequest{
		Input: text,
		Model: "text-embedding-ada-002",
	}
	data, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var embeddingResp EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embeddingResp); err != nil {
		return nil, err
	}

	if len(embeddingResp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return embeddingResp.Data[0].Embedding, nil
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

// Simple whitespace tokenizer
func Tokenize(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

// ComputeF1Score calculates precision, recall, and F1 score
func ComputeF1Score(prediction, reference string) (precision, recall, f1 float64) {
	predTokens := Tokenize(prediction)
	refTokens := Tokenize(reference)

	// Count overlaps
	overlap := 0
	refCounts := make(map[string]int)
	for _, tok := range refTokens {
		refCounts[tok]++
	}
	for _, tok := range predTokens {
		if refCounts[tok] > 0 {
			overlap++
			refCounts[tok]--
		}
	}

	precision = float64(overlap) / float64(len(predTokens))
	recall = float64(overlap) / float64(len(refTokens))
	if precision+recall == 0 {
		f1 = 0
	} else {
		f1 = 2 * precision * recall / (precision + recall)
	}
	return
}

// EvaluateDisharmonySimilarity compares ground truth and GPT response
func EvaluateDisharmonySimilarity(groundTruth, gptOutput string) (float64, error) {
	gtEmb, err := GetEmbedding(groundTruth)
	if err != nil {
		return 0.0, err
	}

	gptEmb, err := GetEmbedding(gptOutput)
	if err != nil {
		return 0.0, err
	}

	similarity := CosineSimilarity(gtEmb, gptEmb)
	return similarity, nil
}

// BatchEvaluate compares GPT responses for all test cases
func BatchEvaluate(testCaseFile string, filename string) error {
	testCases, err := pkg.LoadTestCases(testCaseFile)
	if err != nil {
		return err
	}

	var output bytes.Buffer

	for _, testCase := range testCases {
		// Generate GPT output for each test case
		// -------------------- ZERO SHOT --------------------
		// var fullRegulationText string
		// for _, reg := range testCase.Regulations {
		// 	fullRegulationText += fmt.Sprintf("Document: %s\nArticle: %s\nContent: %s\n\n", reg.Document, reg.Article, reg.Content)
		// }

		// ** --------- No COSTAR ----------
		// var promptBuilder strings.Builder
		// promptBuilder.WriteString("Answer in Indonesian, identify the disharmony between the following law regulations.\n")
		// promptBuilder.WriteString("Input Regulations:\n")
		// promptBuilder.WriteString(fullRegulationText)
		// gptOutput, err := services.EvaluateOpenAIDisharmonyAnalysis(promptBuilder.String())

		// ** --------- COSTAR -----------
		// gptOutput, err := services.EvaluateOpenAIDisharmonyAnalysis(pkg.CostarPromptZeroShot(fullRegulationText))

		// ** -------------------- FEW SHOT --------------------
		// var promptBuilder strings.Builder
		// promptBuilder.WriteString(pkg.CostarHeader)
		// promptBuilder.WriteString("\n\nBelow are some examples:")
		// // Add other test cases as few-shot examples (excluding current)
		// for j, ex := range testCases {
		// 	if i == j {
		// 		continue
		// 	}
		// 	var exampleText strings.Builder
		// 	for _, reg := range ex.Regulations {
		// 		exampleText.WriteString(fmt.Sprintf("Document: %s\nArticle: %s\nContent: %s\n\n", reg.Document, reg.Article, reg.Content))
		// 	}
		// 	promptBuilder.WriteString(fmt.Sprintf(
		// 		"\nInput Regulations:\n%sExpected Disharmony:\n%s\n\n---\n",
		// 		exampleText.String(), ex.Disharmony,
		// 	))
		// }
		// // Add the current test case (as the new input, without disharmony)
		// var currentRegText strings.Builder
		// for _, reg := range testCase.Regulations {
		// 	currentRegText.WriteString(fmt.Sprintf("Document: %s\nArticle: %s\nContent: %s\n\n", reg.Document, reg.Article, reg.Content))
		// }
		// promptBuilder.WriteString("\nNow, analyze the following:\n")
		// promptBuilder.WriteString("Input Regulations:\n")
		// promptBuilder.WriteString(currentRegText.String())
		// promptBuilder.WriteString("End of input.\n")
		// gptOutput, err := services.EvaluateOpenAIDisharmonyAnalysis(promptBuilder.String())

		// ** ----------------------- CHAIN OF THOUGHT -------------------
		// var fullRegulationText string
		// for _, reg := range testCase.Regulations {
		// 	fullRegulationText += fmt.Sprintf("Document: %s\nArticle: %s\nContent: %s\n\n", reg.Document, reg.Article, reg.Content)
		// }
		// gptOutput, err := services.EvaluateOpenAIDisharmonyAnalysis(pkg.CostarPromptChainOfThought(fullRegulationText))

		// ** ----------------------- FEW SHOT + CHAIN OF THOUGHT -------------------
		var promptBuilder strings.Builder
		promptBuilder.WriteString(pkg.CostarHeader)
		promptBuilder.WriteString(`
		Now, analyze the following step by step of reasoning:
		1. Identify the main legal norms, obligations, or permissions from each regulation, including relevant articles or clauses.
		2. Compare these legal norms side by side in terms of:
		   - Definition or terminology
		   - Legal scope and applicability
		   - Authority/responsibility assigned
		   - Exceptions or special conditions
		3. Analyze potential disharmony:
		   - Is there a direct contradiction?
		   - Is there overlapping jurisdiction or authority?
		   - Are there ambiguities that can lead to multiple interpretations?
		4. Conclude with a summary of your findings:
		   - Clearly state if disharmony exists or not.
		   - If disharmony exists, briefly explain the legal or practical implication.
		5. Recommend possible resolutions if appropriate:
		   - Suggest harmonization steps (e.g., revision, repeal, new regulation, clarification).
		   - Mention which regulation may take precedence (if applicable).

		Important Notes:
		1. Be neutral and objective.
		2. Not all case will have disharmony, so if you find no disharmony, please state that clearly.
		3. Provide your analysis in Indonesian, plain text, in paragraph, without numbering, bullets, or any other formatting.
		`)
		// Add other test cases as few-shot examples (excluding current)
		promptBuilder.WriteString("\n\nBelow are some examples:")
		for _, ex := range testCases {
			if ex.ID == testCase.ID {
				continue
			}
			var exampleText strings.Builder
			for _, reg := range ex.Regulations {
				exampleText.WriteString(fmt.Sprintf("Document: %s\nArticle: %s\nContent: %s\n\n", reg.Document, reg.Article, reg.Content))
			}
			promptBuilder.WriteString(fmt.Sprintf(
				"\nInput Regulations:\n%sExpected Disharmony:\n%s\n\n---\n",
				exampleText.String(), ex.Disharmony,
			))
		}
		// Add the current regulations for analysis
		var currentRegText strings.Builder
		for _, reg := range testCase.Regulations {
			currentRegText.WriteString(fmt.Sprintf("Document: %s\nArticle: %s\nContent: %s\n\n", reg.Document, reg.Article, reg.Content))
		}
		promptBuilder.WriteString("\nNow, analyze the regulations input below:\n")
		promptBuilder.WriteString(currentRegText.String())
		promptBuilder.WriteString("End of input.\n")
		gptOutput, err := services.EvaluateOpenAIDisharmonyAnalysis(promptBuilder.String())

		// ------------------------ END OF PROMPT ------------------------
		if err != nil {
			msg := fmt.Sprintf("Error generating output for test case %s: %v\n", testCase.ID, err)
			fmt.Print(msg)
			output.WriteString(msg)
			continue
		}

		// Compare GPT output with ground truth disharmony
		similarity, err := EvaluateDisharmonySimilarity(testCase.Disharmony, gptOutput)
		if err != nil {
			msg := fmt.Sprintf("Error comparing test case %s: %v\n", testCase.ID, err)
			fmt.Print(msg)
			output.WriteString(msg)
			continue
		}

		p, r, f1 := ComputeF1Score(gptOutput, testCase.Disharmony)

		// Build result
		result := fmt.Sprintf(
			"--------------------------------------------------\n"+
				"Test Case: %s - %s\n"+
				"Ground Truth Disharmony:\n%s\n\n"+
				"GPT Output:\n%s\n\n"+
				"Similarity Score: %.4f\n"+
				"Precision: %.4f\n"+
				"Recall: %.4f\n"+
				"F1 Score: %.4f\n\n",
			testCase.ID, testCase.Title, testCase.Disharmony, gptOutput, similarity, p, r, f1,
		)

		if similarity >= 0.80 {
			result += "GPT response semantically matches expected output.\n"
		} else {
			result += "GPT response is semantically different.\n"
		}
		result += ("--------------------------------------------------\n\n")

		output.WriteString(result)
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
	err = BatchEvaluate("../data/test_case.json", outputFile)
	if err != nil {
		fmt.Println("Error during evaluation:", err)
	}
}

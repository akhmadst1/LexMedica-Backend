package main

import (
	"bytes"
	"encoding/json"
	"errors"
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

type ChatDisharmony struct {
	Result   bool   `json:"result"`
	Analysis string `json:"analysis"`
}

// GetEmbedding fetches embedding vector from OpenAI API
func GetEmbedding(text string) ([]float64, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	url := "https://api.openai.com/v1/embeddings"

	payload := EmbeddingRequest{
		Input: text,
		Model: "text-embedding-3-small",
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

func ExtractChatDisharmony(gptOutput string) (ChatDisharmony, error) {
	// Remove any code block formatting like ```json or ``` from GPT output
	cleaned := strings.TrimSpace(gptOutput)
	if strings.HasPrefix(cleaned, "```json") {
		cleaned = strings.TrimPrefix(cleaned, "```json")
	} else if strings.HasPrefix(cleaned, "```") {
		cleaned = strings.TrimPrefix(cleaned, "```")
	}
	cleaned = strings.TrimSuffix(cleaned, "```")
	cleaned = strings.TrimSpace(cleaned)

	var parsed ChatDisharmony
	if err := json.Unmarshal([]byte(cleaned), &parsed); err != nil {
		return ChatDisharmony{}, errors.New("failed to parse JSON from GPT output")
	}
	return parsed, nil
}

// BatchEvaluate compares GPT responses for all test cases
func BatchEvaluate(testCaseFile string, filename string) error {
	testCases, err := pkg.LoadTestCases(testCaseFile)
	if err != nil {
		return err
	}

	var output bytes.Buffer

	var totalSimilarity, totalPrecision, totalRecall, totalF1 float64
	var successfulCases int

	// Generate GPT output for each test case
	for _, testCase := range testCases {
		var fullRegulationText string
		for _, reg := range testCase.Regulations {
			fullRegulationText += fmt.Sprintf("Document: %s\nArticle: %s\nContent: %s\n\n", reg.Document, reg.Article, reg.Content)
		}
		var promptBuilder strings.Builder

		// -------------------- ZERO SHOT --------------------
		promptBuilder.WriteString(pkg.ZeroShot(fullRegulationText))

		// ** -------------------- FEW SHOT --------------------
		// promptBuilder.WriteString(pkg.FewShot(fullRegulationText, testCase.ID))

		// ** ----------------------- CHAIN OF THOUGHT -------------------
		// promptBuilder.WriteString(pkg.ChainOfThought(fullRegulationText))

		// ** ----------------------- FEW SHOT + CHAIN OF THOUGHT -------------------
		// promptBuilder.WriteString(pkg.FewShotChainOfThought(fullRegulationText, testCase.ID))

		// ------------------------ END OF PROMPT ------------------------
		gptOutput, err := services.EvaluateOpenAIDisharmonyAnalysis(promptBuilder.String())
		if err != nil {
			msg := fmt.Sprintf("Error generating GPT output for test case %s: %v\n", testCase.ID, err)
			fmt.Print(msg)
			output.WriteString(msg)
			continue
		}

		parsed, err := ExtractChatDisharmony(gptOutput)
		if err != nil {
			msg := fmt.Sprintf("Error parsing JSON response for test case %s: %v\n", testCase.ID, err)
			fmt.Print(msg)
			output.WriteString(msg)
			continue
		}

		// Compare GPT output with ground truth disharmony
		similarity, err := EvaluateDisharmonySimilarity(testCase.Disharmony, parsed.Analysis)
		if err != nil {
			msg := fmt.Sprintf("Error comparing test case %s: %v\n", testCase.ID, err)
			fmt.Print(msg)
			output.WriteString(msg)
			continue
		}

		p, r, f1 := ComputeF1Score(gptOutput, testCase.Disharmony)

		totalSimilarity += similarity
		totalPrecision += p
		totalRecall += r
		totalF1 += f1
		successfulCases++

		// Build result
		result := fmt.Sprintf(
			"--------------------------------------------------\n"+
				"Test Case: %s - %s\n"+
				"Ground Truth Disharmony:\n%s\n\n"+
				"Disharmony Result:%t\n"+
				"GPT Disharmony Analysis:\n%s\n\n"+
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
				"Average Similarity Score: %.4f\n"+
				"Average Precision:        %.4f\n"+
				"Average Recall:           %.4f\n"+
				"Average F1 Score:         %.4f\n"+
				"=================================================================\n",
			successfulCases, avgSimilarity, avgPrecision, avgRecall, avgF1,
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
	err = BatchEvaluate("../data/test_case.json", outputFile)
	if err != nil {
		fmt.Println("Error during evaluation:", err)
	}
}

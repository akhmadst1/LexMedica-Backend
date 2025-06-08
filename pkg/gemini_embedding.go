package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Gemini Embedding API Structures
type GeminiEmbeddingRequestPart struct {
	Text string `json:"text"`
}

type GeminiEmbeddingRequestContent struct {
	Parts []GeminiEmbeddingRequestPart `json:"parts"`
	Role  string                       `json:"role,omitempty"` // Role is optional for embedding requests
}

type GeminiEmbeddingAPIRequest struct {
	Content GeminiEmbeddingRequestContent `json:"content"`
	// Model string `json:"model"` // Model is in the URL for Gemini, not in the body usually for embedContent
}

type GeminiEmbeddingValue struct {
	Values []float64 `json:"values"` // Or []float32 if you prefer to be more precise
}

type GeminiEmbeddingAPIResponse struct {
	Embedding GeminiEmbeddingValue `json:"embedding"`
}

// GetGeminiEmbedding fetches embedding vector from Gemini API
func GetGeminiEmbedding(text string) ([]float64, error) {
	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}

	// Recommended model: text-embedding-004 (also aliased as embedding-001)
	// You can also use other embedding models if available and suitable.
	modelName := "text-embedding-004"
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:embedContent?key=%s", modelName, geminiKey)

	// Prepare Gemini API payload for embedding
	payload := GeminiEmbeddingAPIRequest{
		Content: GeminiEmbeddingRequestContent{
			Parts: []GeminiEmbeddingRequestPart{
				{Text: text},
			},
			// Role can be omitted or set to "user"
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request payload: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request to Gemini API: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Check for non-200 status codes (API errors)
	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("Gemini API request failed with status %d: %s. Response: %s", resp.StatusCode, resp.Status, string(body))
		// You could try to unmarshal into a GeminiError struct here if you have one defined
		return nil, fmt.Errorf("%s", errMsg)
	}

	var embeddingResp GeminiEmbeddingAPIResponse
	if err := json.Unmarshal(body, &embeddingResp); err != nil {
		return nil, fmt.Errorf("error unmarshalling Gemini embedding response: %w. Raw response: %s", err, string(body))
	}

	if len(embeddingResp.Embedding.Values) == 0 {
		return nil, fmt.Errorf("no embedding vector returned or embedding values are nil/empty. Raw response: %s", string(body))
	}

	return embeddingResp.Embedding.Values, nil
}

package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// GPTEmbeddingRequest is the payload for OpenAI embeddings API
type GPTEmbeddingRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

// GPTEmbeddingResponse holds the API response
type GPTEmbeddingResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
}

// GetOpenAIEmbedding fetches embedding vector from OpenAI API
func GetOpenAIEmbedding(text string) ([]float64, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	url := "https://api.openai.com/v1/embeddings"

	payload := GPTEmbeddingRequest{
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

	var embeddingResp GPTEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embeddingResp); err != nil {
		return nil, err
	}

	if len(embeddingResp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return embeddingResp.Data[0].Embedding, nil
}

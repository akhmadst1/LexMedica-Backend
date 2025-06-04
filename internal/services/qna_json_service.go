package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type QnaRequest struct {
	Question  string   `json:"query"`
	ModelURL  string   `json:"model_url"`
	Embedding string   `json:"embedding_model"`
	History   []string `json:"previous_responses"`
}

func QNAService(req QnaRequest, apiKey string) ([]byte, error) {
	payload, err := json.Marshal(map[string]interface{}{
		"query":              req.Question,
		"embedding_model":    req.Embedding,
		"previous_responses": req.History,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	modelURL := req.ModelURL
	if modelURL == "" {
		modelURL = "http://localhost:8080"
	}

	qnaReq, err := http.NewRequest("POST", modelURL+"/api/chat", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create QnA request: %w", err)
	}
	qnaReq.Header.Set("Content-Type", "application/json")
	qnaReq.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(qnaReq)
	if err != nil {
		return nil, fmt.Errorf("QnA request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("QnA service responded with status: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

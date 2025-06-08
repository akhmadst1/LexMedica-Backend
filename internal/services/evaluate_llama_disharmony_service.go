package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func EvaluateLlamaDisharmonyAnalysis(prompt string) (string, error) {
	ollamaURL := "http://localhost:11434/api/chat"

	// Add `"stream": false` to disable streaming
	payload := map[string]interface{}{
		"model":  "llama3",
		"stream": false,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt + " GIVE YOUR ANSWER IN INDONESIAN LANGUAGE ONLY!",
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", ollamaURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read full response
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var parsed struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}

	if err := json.Unmarshal(rawBody, &parsed); err != nil {
		return "", err
	}

	if parsed.Message.Content != "" {
		return parsed.Message.Content, nil
	}

	return "", fmt.Errorf("no response from LLaMA 3")
}

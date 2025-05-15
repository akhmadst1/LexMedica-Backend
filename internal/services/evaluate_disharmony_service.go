package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func EvaluateOpenAIDisharmonyAnalysis(prompt string) (string, error) {
	openaiKey := os.Getenv("OPENAI_API_KEY")
	openaiUrl := "https://api.openai.com/v1/chat/completions"
	modelName := "gpt-4o-mini"

	payload := []byte(fmt.Sprintf(`{
		"model": "%s",
		"stream": false,
		"messages": [{"role": "user", "content": %q}]
	}`, modelName, prompt))

	req, err := http.NewRequest("POST", openaiUrl, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}

	if len(parsed.Choices) > 0 {
		return parsed.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from GPT")
}

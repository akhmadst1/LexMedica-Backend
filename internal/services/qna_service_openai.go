package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func StreamOpenAIAnswer(question string, w http.ResponseWriter) error {
	openaiKey := os.Getenv("OPENAI_API_KEY")
	openaiUrl := "https://api.openai.com/v1/chat/completions"
	modelName := "gpt-4"

	payload := []byte(fmt.Sprintf(`{
		"model": "%s",
		"stream": true,
		"messages": [{"role": "user", "content": %q}]
	}`, modelName, question))

	req, err := http.NewRequest("POST", openaiUrl, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiKey)

	// Set streaming headers before doing the request
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading stream:", err)
			break
		}

		line = strings.TrimSpace(line)

		if line == "" || line == "data: [DONE]" {
			continue
		}

		if strings.HasPrefix(line, "data: ") {
			jsonPart := strings.TrimPrefix(line, "data: ")
			var parsed struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}
			if err := json.Unmarshal([]byte(jsonPart), &parsed); err != nil {
				fmt.Println("Error parsing chunk:", err)
				continue
			}

			// Send just the content text
			if len(parsed.Choices) > 0 {
				content := parsed.Choices[0].Delta.Content
				if content != "" {
					fmt.Fprintf(w, "data: %s\n\n", content)
					flusher.Flush()
				}
			}
		}
	}

	flusher.Flush()
	return nil
}

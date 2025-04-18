package services

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func StreamLlamaAnswer(question string, w http.ResponseWriter) error {
	ollamaUrl := os.Getenv("OLLAMA_API_URL")
	modelName := "llama3"

	payload := []byte(fmt.Sprintf(`{
		"model": "%s",
		"stream": true,
		"messages": [{"role": "user", "content": "%s"}]
	}`, modelName, question))

	req, err := http.NewRequest("POST", ollamaUrl, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:  true,
			DisableCompression: true,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Critical headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Dev only
	w.Header().Set("Transfer-Encoding", "chunked")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}

	// Stream line by line
	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading stream:", err)
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Only forward JSON lines from Ollama
		if strings.HasPrefix(line, "{") {
			// SSE format
			fmt.Fprintf(w, "data: %s\n\n", line)
			flusher.Flush()
		}
	}

	// End stream
	flusher.Flush()
	return nil
}

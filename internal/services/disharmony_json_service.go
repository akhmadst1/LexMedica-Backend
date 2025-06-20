package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/akhmadst1/tugas-akhir-backend/pkg"
)

func OpenAIDisharmonyAnalysisJSON(regulations string, w http.ResponseWriter) error {
	prompt := pkg.ZeroShot(regulations)
	openaiKey := os.Getenv("OPENAI_API_KEY")
	openaiUrl := "https://api.openai.com/v1/chat/completions"
	modelName := "gpt-4"

	payload := map[string]interface{}{
		"model": modelName,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", openaiUrl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiKey)

	startTime := time.Now()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	durationMs := time.Since(startTime).Milliseconds()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var responseMap struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &responseMap); err != nil {
		return err
	}

	if len(responseMap.Choices) == 0 {
		http.Error(w, "No choices returned", http.StatusInternalServerError)
		return nil
	}

	// Get the assistant's message content
	content := strings.TrimSpace(responseMap.Choices[0].Message.Content)

	// Extract JSON block using regex (in case there's extra explanation text)
	jsonBlockRegex := regexp.MustCompile(`(?s)\{.*\}`)
	jsonMatch := jsonBlockRegex.FindString(content)
	if jsonMatch == "" {
		http.Error(w, "No valid JSON block found in LLM response", http.StatusInternalServerError)
		return nil
	}

	var extracted struct {
		Result   bool   `json:"result"`
		Analysis string `json:"analysis"`
	}
	if err := json.Unmarshal([]byte(jsonMatch), &extracted); err != nil {
		http.Error(w, "Failed to parse JSON from LLM response: "+err.Error(), http.StatusInternalServerError)
		return nil
	}

	// Final response to frontend
	finalResponse := map[string]interface{}{
		"result":             extracted.Result,
		"analysis":           extracted.Analysis,
		"processing_time_ms": durationMs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	return json.NewEncoder(w).Encode(finalResponse)
}

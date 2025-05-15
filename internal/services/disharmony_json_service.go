package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/akhmadst1/tugas-akhir-backend/pkg"
)

func OpenAIDisharmonyAnalysisJSON(regulations string, w http.ResponseWriter) error {
	prompt := pkg.FewShot(regulations)
	openaiKey := os.Getenv("OPENAI_API_KEY")
	openaiUrl := "https://api.openai.com/v1/chat/completions"
	modelName := "gpt-4o-mini"

	// Prepare JSON payload
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

	// Create request
	req, err := http.NewRequest("POST", openaiUrl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read and relay the full JSON response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Write the full response
	w.Write(body)
	return nil
}

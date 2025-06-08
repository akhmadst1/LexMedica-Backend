package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// GeminiRequestPart represents a part of the content in a Gemini request.
type GeminiRequestPart struct {
	Text string `json:"text"`
}

// GeminiRequestContent represents the content to be sent to Gemini.
type GeminiRequestContent struct {
	Role  string              `json:"role,omitempty"`
	Parts []GeminiRequestPart `json:"parts"`
}

// GeminiRequestPayload is the top-level structure for a Gemini API request.
type GeminiRequestPayload struct {
	Contents []GeminiRequestContent `json:"contents"`
	// GenerationConfig *GeminiGenerationConfig `json:"generationConfig,omitempty"` // Optional
	// SafetySettings   []GeminiSafetySetting   `json:"safetySettings,omitempty"`   // Optional
}

// --- Structs for Parsing Gemini Response ---

// GeminiResponsePart corresponds to the "parts" in the Gemini response.
type GeminiResponsePart struct {
	Text string `json:"text"`
	// Add other fields if expecting different part types (e.g., functionCall, fileData)
}

// GeminiResponseContent corresponds to the "content" in a Gemini candidate.
type GeminiResponseContent struct {
	Parts []GeminiResponsePart `json:"parts"`
	Role  string               `json:"role"`
}

// GeminiCandidate corresponds to one of the "candidates" in the Gemini response.
type GeminiCandidate struct {
	Content       GeminiResponseContent `json:"content"`
	FinishReason  string                `json:"finishReason"`
	Index         int                   `json:"index"`
	SafetyRatings []struct {            // Simplified safety ratings
		Category    string `json:"category"`
		Probability string `json:"probability"`
		// Blocked     bool   `json:"blocked,omitempty"` // Present if blocked
	} `json:"safetyRatings"`
}

// GeminiResponse is the top-level structure for a Gemini API response.
type GeminiResponse struct {
	Candidates     []GeminiCandidate `json:"candidates"`
	PromptFeedback *struct {         // Optional
		BlockReason   string `json:"blockReason,omitempty"`
		SafetyRatings []struct {
			Category    string `json:"category"`
			Probability string `json:"probability"`
		} `json:"safetyRatings"`
	} `json:"promptFeedback,omitempty"`
}

// GeminiErrorResponse helps parse error messages from the Gemini API.
type GeminiErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
type GeminiErrorResponse struct {
	Error GeminiErrorDetail `json:"error"`
}

func EvaluateGeminiDisharmonyAnalysis(prompt string) (string, error) {
	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	// Choose a Gemini model.
	// Popular choices: "gemini-1.5-flash-latest", "gemini-1.5-pro-latest", "gemini-pro"
	modelName := "gemini-1.5-flash"
	geminiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent", modelName)

	// Prepare JSON payload for Gemini
	payload := GeminiRequestPayload{
		Contents: []GeminiRequestContent{
			{
				Role: "user", // Specifying the role
				Parts: []GeminiRequestPart{
					{Text: prompt},
				},
			},
		},
		// You can add GenerationConfig here if needed, e.g.:
		// GenerationConfig: &GeminiGenerationConfig{
		//  Temperature: 0.7,
		//  MaxOutputTokens: 256,
		// },
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Gemini payload: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", geminiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", geminiKey) // Gemini API key header

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Gemini: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read Gemini response body: %w", err)
	}

	// Check for API errors indicated by status code or in the body
	if resp.StatusCode >= 400 {
		var apiError GeminiErrorResponse
		if json.Unmarshal(body, &apiError) == nil && apiError.Error.Message != "" {
			return "", fmt.Errorf("error Gemini API (Status %d %s): %s",
				apiError.Error.Code, apiError.Error.Status, apiError.Error.Message)
		}
		return "", fmt.Errorf("request failed Gemini API with status %s. Response: %s", resp.Status, string(body))
	}

	// Parse the successful response
	var parsedResponse GeminiResponse
	if err := json.Unmarshal(body, &parsedResponse); err != nil {
		return "", fmt.Errorf("failed to parse Gemini JSON response: %w. Response body: %s", err, string(body))
	}

	// Check if the prompt or response was blocked due to safety settings
	if parsedResponse.PromptFeedback != nil && parsedResponse.PromptFeedback.BlockReason != "" {
		return "", fmt.Errorf("request Gemini blocked due to prompt safety: %s", parsedResponse.PromptFeedback.BlockReason)
	}

	if len(parsedResponse.Candidates) == 0 {
		// Check if any candidate was blocked
		// The Gemini API might return an empty candidates array if content is blocked,
		// or it might be in PromptFeedback. Some models might have candidate-level finishReason like "SAFETY".
		// It's a bit nuanced. For simplicity, we check common block reasons.
		if parsedResponse.PromptFeedback != nil {
			for _, sr := range parsedResponse.PromptFeedback.SafetyRatings {
				if sr.Probability != "NEGLIGIBLE" && sr.Probability != "LOW" { // Simplified check
					return "", fmt.Errorf("content Gemini generation issue: Prompt safety rating %s for category %s", sr.Probability, sr.Category)
				}
			}
		}
		return "", fmt.Errorf("no candidates found in Gemini response. Response: %s", string(body))
	}

	candidate := parsedResponse.Candidates[0]

	if candidate.FinishReason == "SAFETY" {
		// Log safety ratings for more details if needed
		var safetyIssues []string
		for _, sr := range candidate.SafetyRatings {
			// The API docs usually state "HARM_CATEGORY_..." and "HARM_BLOCK_THRESHOLD_..."
			// Probability is more common in recent API versions: NEGLIGIBLE, LOW, MEDIUM, HIGH
			if sr.Probability != "NEGLIGIBLE" && sr.Probability != "LOW" {
				safetyIssues = append(safetyIssues, fmt.Sprintf("%s: %s", sr.Category, sr.Probability))
			}
		}
		return "", fmt.Errorf("content Gemini generation blocked due to safety. Issues: %v", safetyIssues)
	}

	if candidate.FinishReason != "STOP" && candidate.FinishReason != "MAX_TOKENS" {
		// Other finish reasons might be "OTHER", "UNKNOWN", "RECITATION"
		// "MAX_TOKENS" is acceptable, means it generated up to the limit.
		// "STOP" is the ideal normal completion.
		// You might want to log or handle other reasons specifically.
		// For now, we'll proceed if there's content.
		fmt.Printf("Gemini candidate finished with reason: %s\n", candidate.FinishReason)
	}

	if len(candidate.Content.Parts) == 0 {
		return "", fmt.Errorf("no parts found in Gemini response candidate. Finish reason: %s. Response: %s", candidate.FinishReason, string(body))
	}

	// Assuming the first part contains the text we want.
	// For simple text-in, text-out, this is usually the case.
	// If you expect multiple parts or different types of parts, you'll need to iterate or check types.
	generatedText := candidate.Content.Parts[0].Text

	return generatedText, nil
}

package services

import (
	"github.com/akhmadst1/tugas-akhir-backend/internal/models"
)

// GetLLMAnalysis calls the Disharmony LLM microservice API
// func GetDisharmonyAnalysis(text string) (models.DisharmonyResponse, error) {
// 	client := resty.New()

// 	var response models.DisharmonyResponse
// 	_, err := client.R().
// 		SetHeader("Content-Type", "application/json").
// 		SetBody(models.DisharmonyRequest{Text: text}).
// 		SetResult(&response).
// 		Post("http://disharmony-service:8082/analyze") // Assuming service runs on port 8082

// 	if err != nil {
// 		log.Println("Error calling disharmony service:", err)
// 		return models.DisharmonyResponse{}, err
// 	}

// 	return response, nil
// }

// GetLLMAnalysis returns a dummy analysis result for now
func GetDisharmonyAnalysis(text string) (models.DisharmonyResponse, error) {
	// Mock response
	response := models.DisharmonyResponse{
		Result: "Dummy analysis for the provided regulation text.",
	}
	return response, nil
}

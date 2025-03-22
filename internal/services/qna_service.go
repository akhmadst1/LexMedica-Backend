package services

import (
	"github.com/akhmadst1/tugas-akhir-backend/internal/models"
)

// GetQnAAnswer calls the QnA microservice API
// func GetQnAAnswer(question string) (models.QnAResponse, error) {
// 	client := resty.New()

// 	var response models.QnAResponse
// 	_, err := client.R().
// 		SetHeader("Content-Type", "application/json").
// 		SetBody(models.QnARequest{Question: question}).
// 		SetResult(&response).
// 		Post("http://qna-service:8081/qna") // Assuming service runs on port 8081

// 	if err != nil {
// 		log.Println("Error calling QnA service:", err)
// 		return models.QnAResponse{}, err
// 	}

// 	return response, nil
// }

// GetQnAAnswer returns a dummy response for now
func GetQnAAnswer(question string) (models.QnAResponse, error) {
	// Mock response
	response := models.QnAResponse{
		Answer: "This is a dummy response for: " + question,
	}
	return response, nil
}

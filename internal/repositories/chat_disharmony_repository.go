package repositories

import (
	"fmt"

	"github.com/akhmadst1/tugas-akhir-backend/config"
	"github.com/akhmadst1/tugas-akhir-backend/internal/models"
)

func CreateChatDisharmony(messageID int, result bool, analysis string, processingTimeMs int) (models.DisharmonyAnalysis, error) {
	var insertedDisharmonyAnalysis []models.DisharmonyAnalysis

	err := config.Supabase.DB.
		From("disharmony_analysis").
		Insert(map[string]interface{}{
			"message_id":       messageID,
			"result":           result,
			"analysis":         analysis,
			"processingTimeMs": processingTimeMs,
		}).
		Execute(&insertedDisharmonyAnalysis)

	if err != nil || len(insertedDisharmonyAnalysis) == 0 {
		return models.DisharmonyAnalysis{}, fmt.Errorf("failed to insert disharmony analysis to chat: %w", err)
	}

	return insertedDisharmonyAnalysis[0], nil
}

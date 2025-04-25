package models

type DisharmonyAnalysisRequest struct {
	Regulations string `json:"regulations"` // Accepts any number of regs
	Method      string `json:"method"`      // Prompt method
}

type DisharmonyAnalysisResponse struct {
	Result string `json:"result"`
}

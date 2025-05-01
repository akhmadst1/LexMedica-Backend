package models

type DisharmonyAnalysisRequest struct {
	Method      string `json:"method"`      // Prompt method
	Regulations string `json:"regulations"` // Accepts any number of regs
}

type DisharmonyAnalysisResponse struct {
	Result string `json:"result"`
}

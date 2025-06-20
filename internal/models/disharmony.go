package models

type DisharmonyAnalysisRequest struct {
	Regulations string `json:"regulations"` // Accepts any number of regs
}

type DisharmonyAnalysisResponse struct {
	Result string `json:"result"`
}

package dto


// GenerateRequest represents the request body for content generation.
type GenerateRequest struct {
	Prompt      string  `json:"prompt" binding:"required"`
	VoiceDNA    string  `json:"voice_dna" binding:"required"`
}

// GenerateResponse represents the response body for generated content.
type GenerateResponse struct {
	Content string `json:"content"`
	Status  string `json:"status"`
}
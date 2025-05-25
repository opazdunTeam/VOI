package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TranscribeRequest struct {
	AudioData []byte `json:"audio_data" binding:"required"`
	Language  string `json:"language"`
}

type TranscribeResponse struct {
	Text   string `json:"text"`
	Status string `json:"status"`
}

func HandleTranscribe(c *gin.Context) {
	var req TranscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement transcription logic using service layer

	response := TranscribeResponse{
		Text:   "Transcribed text will appear here",
		Status: "success",
	}

	c.JSON(http.StatusOK, response)
}

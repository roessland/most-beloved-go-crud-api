package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateQuoteInput struct {
	Book  string `json:"book" binding:"required"`
	Quote string `json:"quote" binding:"required"`
}

type CreateQuoteResult struct {
	UUID  string `json:"uuid"`
	Book  string `json:"book"`
	Quote string `json:"quote"`
}

func Health(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

// Create creates a new quote
func Create(c *gin.Context) {
	var input CreateQuoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := createQuote(input)

	c.JSON(http.StatusCreated, result)
}

func createQuote(request CreateQuoteInput) CreateQuoteResult {
	// TODO:: Save to DB
	return CreateQuoteResult{
		UUID:  "1234-TODO-1234",
		Book:  request.Book,
		Quote: request.Quote,
	}
}

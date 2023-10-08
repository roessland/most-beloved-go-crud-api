package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/roessland/most-beloved-go-crud-api/db"
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

// Quotes holds resources used by handlers
type Quotes struct {
	Queries *db.Queries
}

// Create creates a new quote
func (quotes *Quotes) Create(c *gin.Context) {
	var input CreateQuoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := quotes.create(c, input)
	if err != nil {
		// TODO: Don't return the actual error in release mode
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// create creates a new quote
func (quotes *Quotes) create(ctx context.Context, input CreateQuoteInput) (*CreateQuoteResult, error) {
	dbParams := db.CreateQuoteParams{
		ID:         uuid.New(),
		Book:       input.Book,
		Quote:      input.Quote,
		InsertedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	dbQuote, err := quotes.Queries.CreateQuote(ctx, dbParams)
	if err != nil {
		return nil, err
	}
	return &CreateQuoteResult{
		UUID:  dbQuote.ID.String(),
		Book:  dbQuote.Book,
		Quote: dbQuote.Quote,
	}, nil
}

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package db

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Quote struct {
	ID         uuid.UUID
	Book       string
	Quote      string
	InsertedAt pgtype.Timestamptz
	UpdatedAt  pgtype.Timestamptz
}
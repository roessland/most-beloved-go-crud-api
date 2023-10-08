
-- name: CreateQuote :one
INSERT INTO quotes (id, book, quote, inserted_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

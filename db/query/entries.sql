-- name: CreateEntries :one
-- Create a new entries
INSERT INTO entries (account_id, amount) VALUES ($1, $2) RETURNING *;

-- name: GetEntries :one
-- Get an entries by id
SELECT * FROM entries WHERE id = $1;

-- name: ListEntries :many
-- List all entries
SELECT * FROM entries
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateEntries :one
UPDATE entries
SET amount = $2
WHERE id = $1
RETURNING *;

-- name: DeleteEntries :exec
-- Delete an entries
DELETE FROM entries WHERE id = $1;
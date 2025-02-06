-- name: CreateTransfers :one
-- Create a new transfers
INSERT INTO transfers (from_account_id, to_account_id, amount) VALUES ($1, $2, $3) RETURNING *;

-- name: GetTransfers :one
-- Get a transfers by id
SELECT * FROM transfers WHERE id = $1;

-- name: ListTransfers :many
SELECT t.*
FROM transfers t
JOIN account a ON t.from_account_id = a.id
WHERE a.owner = $1
ORDER BY t.id
LIMIT $2
OFFSET $3;

-- name: UpdateTransfer :one
UPDATE transfers
SET amount = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTransfers :exec
-- Delete a transfers
DELETE FROM transfers WHERE id = $1;
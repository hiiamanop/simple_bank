-- name: CreateAccount :one
-- Create a new account
INSERT INTO account (owner, balance, currency) VALUES ($1, $2, $3) RETURNING *;

-- name: GetAccount :one
-- Get an account by id
SELECT * FROM account WHERE id = $1;

-- name: ListAccounts :many
-- List all accounts
SELECT * FROM account
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: AddAccountBalance :one
UPDATE account
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
-- Delete an account
DELETE FROM account WHERE id = $1;


-- name: CreateUser :one
-- Create a new User
INSERT INTO users (username, hashed_password, fullname, email) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetUser :one
-- Get an User by id
SELECT * FROM users WHERE username = $1 LIMIT 1;
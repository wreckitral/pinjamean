-- name: GetAccountByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: SaveAccount :exec
INSERT INTO users (uuid, email, passwordHashed, role, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

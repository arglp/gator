-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
Values (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE name = $1;

-- name: GetUsers :many
SELECT * FROM users
ORDER BY created_at DESC;

-- name: DeleteUsers :exec
DELETE FROM users;
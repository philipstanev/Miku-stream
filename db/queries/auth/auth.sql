-- name: InsertUser :one
INSERT INTO users(password_hash, username)
VALUES ($1, $2)
RETURNING id;

-- name: GetUserByID :one
SELECT username, role, password_hash, id FROM users
WHERE id=$1;

-- name: GetUserByUsername :one
SELECT username, role, password_hash, id FROM users
WHERE username=$1;

-- name: SetRole :exec
UPDATE users 
SET role=$1
WHERE id=$2;



-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html#  

-- name: CreateUser :one
INSERT INTO users (
  username, 
  hashed_password, 
  full_name, 
  email
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;


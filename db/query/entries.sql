-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html#  

-- name: CreateEntrie :one
INSERT INTO entries (
  account_id, 
  amount
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetEntrie :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY id
LIMIT $1 
OFFSET $2;

-- name: UpdateAEntrieNoReturn :exec
UPDATE entries
SET account_id =$2, amount = $3
WHERE id = $1;

-- name: UpdateEntrie :one
UPDATE entries
SET account_id =$2, amount = $3
WHERE id = $1
RETURNING *;

-- name: DeleteEntrie :exec
DELETE FROM entries
WHERE id = $1;
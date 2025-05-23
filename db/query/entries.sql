-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html#  

-- name: CreateEntry :one
INSERT INTO entries (
  account_id, 
  amount
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
WHERE account_id = $1
ORDER BY id
LIMIT $2 
OFFSET $3;

-- name: UpdateAEntrieNoReturn :exec
UPDATE entries
SET account_id =$2, amount = $3
WHERE id = $1;

-- name: UpdateEntry :one
UPDATE entries
SET account_id =$2, amount = $3
WHERE id = $1
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;
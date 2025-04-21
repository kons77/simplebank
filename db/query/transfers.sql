-- https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html#  

-- name: CreateTrasfer :one
INSERT INTO transfers (
  from_account_id, 
  to_account_id,
  amount
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetTrasfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTrasfers :many
SELECT * FROM transfers
ORDER BY id
LIMIT $1 
OFFSET $2;

-- name: UpdateATrasferNoReturn :exec
UPDATE transfers
SET from_account_id =$2, to_account_id = $3, amount = $4
WHERE id = $1;

-- name: UpdateTrasfer :one
UPDATE transfers
SET from_account_id =$2, to_account_id = $3, amount = $4
WHERE id = $1
RETURNING *;

-- name: DeleteTrasfer :exec
DELETE FROM transfers
WHERE id = $1;
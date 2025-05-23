// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: entries.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createEntry = `-- name: CreateEntry :one

INSERT INTO entries (
  account_id, 
  amount
) VALUES (
  $1, $2
) RETURNING id, account_id, amount, created_at
`

type CreateEntryParams struct {
	AccountID pgtype.Int8 `json:"account_id"`
	Amount    int64       `json:"amount"`
}

// https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html#
func (q *Queries) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error) {
	row := q.db.QueryRow(ctx, createEntry, arg.AccountID, arg.Amount)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const deleteEntry = `-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1
`

func (q *Queries) DeleteEntry(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteEntry, id)
	return err
}

const getEntry = `-- name: GetEntry :one
SELECT id, account_id, amount, created_at FROM entries
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetEntry(ctx context.Context, id int64) (Entry, error) {
	row := q.db.QueryRow(ctx, getEntry, id)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const listEntries = `-- name: ListEntries :many
SELECT id, account_id, amount, created_at FROM entries
WHERE account_id = $1
ORDER BY id
LIMIT $2 
OFFSET $3
`

type ListEntriesParams struct {
	AccountID pgtype.Int8 `json:"account_id"`
	Limit     int32       `json:"limit"`
	Offset    int32       `json:"offset"`
}

func (q *Queries) ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error) {
	rows, err := q.db.Query(ctx, listEntries, arg.AccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Entry{}
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAEntrieNoReturn = `-- name: UpdateAEntrieNoReturn :exec
UPDATE entries
SET account_id =$2, amount = $3
WHERE id = $1
`

type UpdateAEntrieNoReturnParams struct {
	ID        int64       `json:"id"`
	AccountID pgtype.Int8 `json:"account_id"`
	Amount    int64       `json:"amount"`
}

func (q *Queries) UpdateAEntrieNoReturn(ctx context.Context, arg UpdateAEntrieNoReturnParams) error {
	_, err := q.db.Exec(ctx, updateAEntrieNoReturn, arg.ID, arg.AccountID, arg.Amount)
	return err
}

const updateEntry = `-- name: UpdateEntry :one
UPDATE entries
SET account_id =$2, amount = $3
WHERE id = $1
RETURNING id, account_id, amount, created_at
`

type UpdateEntryParams struct {
	ID        int64       `json:"id"`
	AccountID pgtype.Int8 `json:"account_id"`
	Amount    int64       `json:"amount"`
}

func (q *Queries) UpdateEntry(ctx context.Context, arg UpdateEntryParams) (Entry, error) {
	row := q.db.QueryRow(ctx, updateEntry, arg.ID, arg.AccountID, arg.Amount)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

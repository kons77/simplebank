package util

import "github.com/jackc/pgx/v5/pgtype"

// toPgInt8 convert int64 to pgtype.Int8
func ToPgInt8(v int64) pgtype.Int8 {
	return pgtype.Int8{
		Int64: v,
		Valid: true,
	}
}

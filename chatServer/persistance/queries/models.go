// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package queries

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Message struct {
	ID       int32
	Content  string
	Senderid pgtype.Int4
}

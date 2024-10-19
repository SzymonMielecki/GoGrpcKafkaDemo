package persistance

import (
	"context"

	"github.com/SzymonMielecki/GoGrpcKafkaDemo/chatServer/persistance/queries"
	"github.com/SzymonMielecki/GoGrpcKafkaDemo/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type DB struct {
	*queries.Queries
}

func NewDB(host, user, password, dbname, port string) (*DB, error) {
	ctx := context.Background()
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable TimeZone=Europe/Warsaw"
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	queries := queries.New(conn)
	return &DB{queries}, nil
}

func (db *DB) CreateMessage(message *types.Message) (*types.Message, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db.Queries.CreateMessage(ctx, queries.CreateMessageParams{
		Content:  message.Content,
		Senderid: pgtype.Int4{Int32: int32(message.SenderID), Valid: true},
	})
	return message, nil
}

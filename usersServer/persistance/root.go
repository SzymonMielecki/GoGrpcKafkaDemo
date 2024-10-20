package persistance

import (
	"context"
	"fmt"

	"github.com/SzymonMielecki/GoGrpcKafkaDemo/types"
	"github.com/SzymonMielecki/GoGrpcKafkaDemo/usersServer/persistance/queries"
	"github.com/jackc/pgx/v5"
)

type DB struct {
	Queries *queries.Queries
	Conn    *pgx.Conn
}

func NewDB(host, user, password, dbname, port string) (*DB, error) {
	ctx := context.Background()
	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	fmt.Println(url)
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	queries := queries.New(conn)
	return &DB{Queries: queries, Conn: conn}, nil
}

func (db *DB) Close() {
	db.Conn.Close(context.Background())
}

func (db *DB) CreateUser(user *types.User) (*types.User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	u, err := db.Queries.CreateUser(ctx, queries.CreateUserParams{
		Username:     user.Username,
		Email:        user.Email,
		Passwordhash: user.PasswordHash,
	})
	if err != nil {
		return nil, err
	}
	return &types.User{
		ID:           uint(u.ID),
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.Passwordhash,
	}, nil
}

func (db *DB) GetUserById(id uint) (*types.User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fmt.Println("ID: ", id)
	u, err := db.Queries.GetUserById(ctx, int32(id))
	if err != nil {
		return nil, err
	}
	return &types.User{
		ID:           uint(u.ID),
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.Passwordhash,
	}, nil
}
func (db *DB) GetUserByEmail(email string) (*types.User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	u, err := db.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &types.User{
		ID:           uint(u.ID),
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.Passwordhash,
	}, nil
}

func (db *DB) GetUserByUsername(username string) (*types.User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	u, err := db.Queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return &types.User{
		ID:           uint(u.ID),
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.Passwordhash,
	}, nil
}

func (db *DB) GetUserByUsernameAndEmail(username, email string) (*types.User, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	u, err := db.Queries.GetUserByUsernameAndEmail(ctx, queries.GetUserByUsernameAndEmailParams{
		Username: username,
		Email:    email,
	})
	if err != nil {
		return nil, err
	}
	return &types.User{
		ID:           uint(u.ID),
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.Passwordhash,
	}, nil
}

func (db *DB) UsernameExists(username string) bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b, err := db.Queries.UsernameExists(ctx, username)
	return err != nil && b
}

func (db *DB) EmailExists(email string) bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b, err := db.Queries.EmailExists(ctx, email)
	return err != nil && b
}

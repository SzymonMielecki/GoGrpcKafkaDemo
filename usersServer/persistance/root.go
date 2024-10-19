package persistance

import (
	"context"

	"github.com/SzymonMielecki/GoGrpcKafkaDemo/types"
	"github.com/SzymonMielecki/GoGrpcKafkaDemo/usersServer/persistance/queries"
	"github.com/jackc/pgx/v5"
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
	return &DB{Queries: queries}, nil
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

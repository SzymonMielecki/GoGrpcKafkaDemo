// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package queries

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO
    Users (Username, Email, PasswordHash)
VALUES ($1, $2, $3) RETURNING id, email, username, passwordhash
`

type CreateUserParams struct {
	Username     string
	Email        string
	Passwordhash string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Username, arg.Email, arg.Passwordhash)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Passwordhash,
	)
	return i, err
}

const emailExists = `-- name: EmailExists :one
SELECT EXISTS ( SELECT 1 FROM Users WHERE Email = $1 )
`

func (q *Queries) EmailExists(ctx context.Context, email string) (bool, error) {
	row := q.db.QueryRow(ctx, emailExists, email)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, username, passwordhash FROM Users WHERE Email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Passwordhash,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, email, username, passwordhash FROM Users WHERE Id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Passwordhash,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, email, username, passwordhash FROM Users WHERE Username = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Passwordhash,
	)
	return i, err
}

const getUserByUsernameAndEmail = `-- name: GetUserByUsernameAndEmail :one
SELECT id, email, username, passwordhash FROM Users WHERE Username = $1 AND Email = $2
`

type GetUserByUsernameAndEmailParams struct {
	Username string
	Email    string
}

func (q *Queries) GetUserByUsernameAndEmail(ctx context.Context, arg GetUserByUsernameAndEmailParams) (User, error) {
	row := q.db.QueryRow(ctx, getUserByUsernameAndEmail, arg.Username, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Passwordhash,
	)
	return i, err
}

const getUserByUsernameAndPassword = `-- name: GetUserByUsernameAndPassword :one
SELECT id, email, username, passwordhash FROM Users WHERE Username = $1 AND PasswordHash = $2
`

type GetUserByUsernameAndPasswordParams struct {
	Username     string
	Passwordhash string
}

func (q *Queries) GetUserByUsernameAndPassword(ctx context.Context, arg GetUserByUsernameAndPasswordParams) (User, error) {
	row := q.db.QueryRow(ctx, getUserByUsernameAndPassword, arg.Username, arg.Passwordhash)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Passwordhash,
	)
	return i, err
}

const usernameExists = `-- name: UsernameExists :one
SELECT EXISTS ( SELECT 1 FROM Users WHERE Username = $1 )
`

func (q *Queries) UsernameExists(ctx context.Context, username string) (bool, error) {
	row := q.db.QueryRow(ctx, usernameExists, username)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

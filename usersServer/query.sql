-- name: CreateUser :one
INSERT INTO
    Users (Username, Email, PasswordHash)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetUserById :one
SELECT * FROM Users WHERE Id = $1;

-- name: GetUserByEmail :one
SELECT * FROM Users WHERE Email = $1;

-- name: GetUserByUsername :one
SELECT * FROM Users WHERE Username = $1;

-- name: GetUserByUsernameAndEmail :one
SELECT * FROM Users WHERE Username = $1 AND Email = $2;

-- name: GetUserByUsernameAndPassword :one
SELECT * FROM Users WHERE Username = $1 AND PasswordHash = $2;

-- name: UsernameExists :one
SELECT EXISTS ( SELECT 1 FROM Users WHERE Username = $1 );

-- name: EmailExists :one
SELECT EXISTS ( SELECT 1 FROM Users WHERE Email = $1 );
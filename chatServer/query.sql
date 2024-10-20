-- name: CreateMessage :exec
INSERT INTO messages (Content, SenderId) VALUES ($1, $2);
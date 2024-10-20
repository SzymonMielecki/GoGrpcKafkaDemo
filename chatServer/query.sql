-- name: CreateMessage :exec
INSERT INTO Messages (Content, SenderId) VALUES ($1, $2);
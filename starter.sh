#!/bin/bash

USERS_INPUT_FILE="usersServer/schema.sql"
USERS_OUTPUT_FILE=".init-scripts/init-users-db.sql"

CHAT_INPUT_FILE="chatServer/schema.sql"
CHAT_OUTPUT_FILE=".init-scripts/init-chat-db.sql"

mkdir -p init-scripts

process_schema() {
    local input_file=$1
    local output_file=$2

    awk '
    {
        if ($0 ~ /CREATE TABLE/) {
            sub(/CREATE TABLE/, "CREATE TABLE IF NOT EXISTS")
        }
        print
    }' "$input_file" > "$output_file"

    echo "Schema has been copied and modified to $output_file"
}

go build -o chatUp client/main.go

process_schema "$USERS_INPUT_FILE" "$USERS_OUTPUT_FILE"
process_schema "$CHAT_INPUT_FILE" "$CHAT_OUTPUT_FILE"

docker-compose up --build
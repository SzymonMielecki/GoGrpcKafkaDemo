#!/bin/sh
set -e

db_host="$1"
db_port="$2"
kafka_host="$3"
kafka_port="$4"
shift 4
cmd="$@"

until nc -z "$db_host" "$db_port"; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

until nc -z "$kafka_host" "$kafka_port"; do
  >&2 echo "Kafka is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres and Kafka are up - executing command"
exec $cmd
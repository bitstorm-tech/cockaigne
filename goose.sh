#! /bin/sh

set -a
. .env
set +a

export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="user=$PG_USER password=$PG_PASSWORD dbname=$PG_DATABASE sslmode=disable"
export GOOSE_MIGRATION_DIR=database/migrations

goose "$@"

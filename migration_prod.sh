#!/bin/bash
source local.env

sleep 2 && goose -dir "${MIGRATION_DIR}" postgres "host=pg-prod port=5432 dbname=$POSTGRES_DB user=$POSTGRES_USER password=$POSTGRES_PASSWORD sslmode=disable" up -v
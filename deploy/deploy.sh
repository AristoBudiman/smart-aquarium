#!/bin/sh
set -e
git pull
docker compose --env-file .env -f deploy/docker-compose.yml build backend
docker compose --env-file .env -f deploy/docker-compose.yml build caddy
docker compose --env-file .env -f deploy/docker-compose.yml up -d
docker image prune -f


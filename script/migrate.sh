#!/usr/bin/env bash

# Load environment variables from .env
if [ -f .env ]; then
    set -a
    source .env
    set +a
fi

command="$1"
name="$2"

case "$command" in
    up)
        migrate -path migrations -database "$DATABASE_URL" up
        ;;
    down)
        count="${name:-1}"
        read -rp "Rolling back $count migration(s). Continue? [y/N] " confirm
        if [[ "$confirm" == "y" || "$confirm" == "Y" ]]; then
            migrate -path migrations -database "$DATABASE_URL" down "$count"
        fi
        ;;
    create)
        if [ -z "$name" ]; then
            echo "Usage: $0 create <migration_name>"
            exit 1
        fi
        migrate create -ext sql -dir migrations -seq "$name"
        ;;
    force)
        if [ -z "$name" ]; then
            echo "Usage: $0 force <version>"
            exit 1
        fi
        migrate -path migrations -database "$DATABASE_URL" force "$name"
        ;;
    *)
        cat <<EOF
Usage: $0 <command> [argument]

Commands:
  up                  Apply all pending migrations
  down [count]        Roll back count migrations (default: 1)
  create <name>       Create a new migration
  force <version>     Force migration version
EOF
        exit 1
        ;;
esac

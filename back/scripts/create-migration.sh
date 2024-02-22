#!/usr/bin/env bash

NAME=""
MIGRATION_CONTENT="BEGIN;

COMMIT;"

usage() {
    cat >&2 <<EOF
Usage: $0 [options]

Creates a new golang-migrate migration.

Supported options:
  --help                    Prints this message and exits.
  -n, --name                The migration name

EOF
}

while [[ $# -ne 0 ]]
do
  case "$1" in
    -n|--name)
      shift
      if [[ -z "$1" ]]; then
        echo "empty argument to --name flag" >&2
        exit 1
      fi
      NAME="$1"
      ;;
    --help)
      usage
      exit
      ;;
    *)
      echo "unknown option: $1" >&2
      usage
      exit 2
      ;;
  esac
  shift
done

if [[ -z "${NAME}" ]]; then
  echo "must specify a --name flag" >&2
  exit 1
fi

EPOCH=$(date +%s)

# ROOT_DIR=$(git rev-parse --show-toplevel || echo ".")
# MIGRATION_DIR="${ROOT_DIR}/internal/sql/migrations"
# ROOT_DIR="."
MIGRATION_DIR="../sql/migrations"

echo "${MIGRATION_CONTENT}" > "${MIGRATION_DIR}/${EPOCH}_${NAME}.down.sql"
echo "${MIGRATION_CONTENT}" > "${MIGRATION_DIR}/${EPOCH}_${NAME}.up.sql"
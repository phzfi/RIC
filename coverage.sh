#!/bin/bash
# Report test coverage
set -e
cd "$(dirname "$0")"

workspace="."
mkdir -p ${workspace}/reports/coverage
rm -rf ${workspace}/reports/coverage/*

docker compose run --rm test

sudo chown -R $(whoami):$(id -gn) ${workspace}/reports/coverage || true

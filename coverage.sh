#!/bin/bash
# Report test coverage
set -e
cd "$(dirname "$0")"

workspace="."
mkdir -p ${workspace}/reports/coverage

docker compose run --rm test sh -c 'go test -v -coverprofile=reports/coverage/coverage.out ./server/... && go tool cover -html=reports/coverage/coverage.out -o=reports/coverage/coverage.html'

sudo chown -R $(whoami):$(id -gn) ${workspace}/reports/coverage || true
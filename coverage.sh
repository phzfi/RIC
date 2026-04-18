#!/bin/bash
# Report test coverage
USER=`whoami`
GROUP=`id -gn`
# Check if workspace is set
if [ -z "$workspace" ]; then
  # If not set, default to current directory
  workspace="."
fi

mkdir -p ${workspace}/reports/coverage
docker compose run --rm test sh -c 'go test -v -coverprofile=reports/coverage/coverage.out ./server/... && go tool cover -html=reports/coverage/coverage.out -o=reports/coverage/coverage.html'
sudo chown -R $USER:$GROUP ${workspace}/reports/coverage || true

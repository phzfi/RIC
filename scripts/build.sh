#!/bin/bash

cd "$GOPATH/src/github.com/phzfi/RIC/server"
go get -t -v ./...
go build

mkdir -p /var/lib/ric/images
mkdir -p /var/lib/ric/cache

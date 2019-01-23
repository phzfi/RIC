#!/bin/bash

export GOPATH=/root/go

echo "GOPATH: $GOPATH"
cd /root/go/src/phzfi/RIC/server
go get -t -v ./...
go build


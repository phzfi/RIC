#!/bin/bash

export GOPATH=/root/go

cd "$GOPATH/src/github.com/phzfi/RIC/server"
go test

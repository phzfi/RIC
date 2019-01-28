#!/bin/bash

export GOPATH=/root/go
export PATH=$PATH:$GOPATH/bin

cd "$GOPATH/src/github.com/phzfi/RIC/server"

go get -u github.com/jstemmer/go-junit-report
go get github.com/t-yuki/gocover-cobertura


go test -v -coverprofile=cover.out ./ 2>&1 |  go-junit-report > junit.xml

gocover-cobertura < cover.out > coverage.xml
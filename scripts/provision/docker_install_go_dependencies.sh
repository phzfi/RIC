#!/bin/bash

cd /root/go/src/github.com/phzfi/RIC/server
go get -t -v ./...
go get -v github.com/derekparker/delve/cmd/dlv
go build
go build -tags debug -v -gcflags "all=-N -l"


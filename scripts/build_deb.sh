#!/bin/bash

export GOPATH=/root/go

cd "$GOPATH/src/github.com/phzfi/RIC/server"
go build

cp "$GOPATH/src/github.com/phzfi/RIC/server/server" "/ric/build/deb/usr/local/bin/ric"
cp /ric/config/host_whitelist.ini /ric/build/deb/etc/ric/host_whitelist.ini
cp /ric/config/ric_config.ini /ric/build/deb/etc/ric/ric_config.ini

cd "/ric/build"

dpkg-deb --build deb  /ric/build/phz-ric.deb

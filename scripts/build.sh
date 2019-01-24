#!/bin/bash

export GOPATH=/root/go

echo "GOPATH: $GOPATH"
cd /root/go/src/github.com/phzfi/RIC/server
go get -t -v ./...
go build

mkdir -p /var/lib/ric/images
mkdir -p /var/lib/ric/cache
mkdir -p /var/lib/ric/config

cp -R /root/go/src/github.com/phzfi/RIC/mount/etc/ric /etc/ric
cp -R /root/go/src/github.com/phzfi/RIC/mount/config /var/lib/ric/config

sed -i 's|/mnt/RIC_image_repository|/var/lib/ric|g' /etc/ric/ric_config.ini
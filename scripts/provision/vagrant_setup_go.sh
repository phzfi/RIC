#!/bin/bash
OPWD="$(pwd)"
apt-get install -y golang-go golang-go.tools golang-golang-x-tools-dev

echo "export GOPATH=/home/vagrant/go" >> /home/vagrant/.bashrc
export GOPATH="/home/vagrant/go"

echo "export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin" >> /home/vagrant/.bashrc
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

mkdir -p "/home/vagrant/go/bin"
mkdir -p "/home/vagrant/go/pkg"

#apt-get install libmagickwand-dev libmagickcore-dev #imagemagick

#apt-get install gccgo-go

cd /home/vagrant/go/src/github.com/phzfi/RIC/
go get -t ./...

go get github.com/derekparker/delve/cmd/dlv

## Ownership
chown -R vagrant:vagrant /home/vagrant/go

cd "${OPWD}"


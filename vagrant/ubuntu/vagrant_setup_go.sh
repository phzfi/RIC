#!/bin/bash
OPWD="$(pwd)"
apt-get install -y golang-go golang-go.tools golang-golang-x-tools-dev

echo "export GOPATH=/home/vagrant/go" >> /home/vagrant/.bashrc
export GOPATH="/home/vagrant/go"
mkdir -p "/home/vagrant/go/bin"
mkdir -p "/home/vagrant/go/pkg"

cd /home/vagrant/go/src/github.com/phzfi/RIC/
go get -t ./...

# Ownership
chown -R vagrant:vagrant /home/vagrant/go

cd "${OPWD}"


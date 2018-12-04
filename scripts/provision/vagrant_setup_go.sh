#!/bin/bash
OPWD="$(pwd)"
#apt-get install -y golang-go golang-go.tools golang-golang-x-tools-dev

mkdir golang_install
cd golang_install
curl -O https://storage.googleapis.com/golang/go1.11.2.linux-amd64.tar.gz
tar -xvf go1.11.2.linux-amd64.tar.gz
sudo mv go /usr/local/
rm go1.11.2.linux-amd64.tar.gz
cd ..
rmdir golang_install


echo "export GOPATH=/home/vagrant/go" >> /home/vagrant/.bashrc
export GOPATH="/home/vagrant/go"

echo "export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin" >> /home/vagrant/.bashrc
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

mkdir -p "/home/vagrant/go/bin"
mkdir -p "/home/vagrant/go/pkg"


cd /home/vagrant/go/src/github.com/phzfi/RIC/
go get -t ./...

go get github.com/derekparker/delve/cmd/dlv

## Ownership
chown -R vagrant:vagrant /home/vagrant/go

cd "${OPWD}"


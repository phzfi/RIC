#!/bin/bash

pacman -Syyu --noconfirm
pacman -S imagemagick go libwebp --noconfirm

echo "export GOPATH=/home/vagrant/go" >> /home/vagrant/.bashrc
export GOPATH=/home/vagrant/go

cd /home/vagrant/go/src/github.com/phzfi/RIC/
go get -t ./...

chown -R vagrant:vagrant /home/vagrant/go

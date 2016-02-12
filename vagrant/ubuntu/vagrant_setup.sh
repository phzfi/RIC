#!/bin/bash

# Install the absolute requirements
apt-get update
apt-get install -y libmagickwand-dev libwebp-dev libtiff5-dev libpng12-dev
apt-get install -y dpkg-dev build-essential devscripts
apt-get install -y golang-go golang-go.tools golang-golang-x-tools-dev

echo "export GOPATH=/home/vagrant/go" >> /home/vagrant/.bashrc
export GOPATH="/home/vagrant/go"
mkdir -p "/home/vagrant/go/bin"
mkdir -p "/home/vagrant/go/pkg"
mkdir -p "/home/vagrant/go/src"

cd /home/vagrant/go/src/github.com/phzfi/RIC/
go get -t ./...

# Ownership
chown -R vagrant:vagrant /home/vagrant/go

mkdir -p /tmp/imagemagick
cd /tmp/imagemagick
apt-get build-dep imagemagick -y
apt-get source imagemagick -y
cd imagemagick-*
debuild -uc -us
dpkg -i ../*magick*.deb



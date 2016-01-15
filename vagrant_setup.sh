#!/bin/bash
sudo apt-get update
sudo apt-get install -y libmagickwand-dev
sudo apt-get install -y golang-go dpkg-dev build-essential
echo "export GOPATH=/home/vagrant/go" >> /home/vagrant/.bashrc
export GOPATH=/home/vagrant/go
chown -R vagrant:vagrant /home/vagrant/go

cd /home/vagrant/go/src/github.com/phzfi/RIC/
go get -t ./...

cd /tmp
mkdir imagemagick
cd imagemagick
sudo apt-get build-dep imagemagick -y
sudo apt-get install libwebp-dev devscripts -y
apt-get source imagemagick -y
cd imagemagick-*
debuild -uc -us
sudo dpkg -i ../*magick*.deb


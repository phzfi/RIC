#!/bin/bash

sudo apt-get -y install wget
wget -O - http://pkg.phz.fi/pkg.phz.fi.gpg.key | sudo apt-key add -
# wget http://pkg.phz.fi/pkg.phz.fi.gpg.key -O pkg.key && sudo apt-key add pkg.key
sudo mkdir -p /etc/apt/sources.list.d/
echo 'echo "deb http://pkg.phz.fi/bionic ./" >> /etc/apt/sources.list.d/pkg.phz.fi.list' | sudo -s
sudo apt-get update
sudo apt-get -y install phz-ric
sudo apt-get -y install nfs-common

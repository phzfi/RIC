#!/bin/bash

sudo apt-get -y install software-properties-common
sudo apt-get -y install wget
wget -O - http://pkg.phz.fi/pkg.phz.fi.gpg.key | sudo apt-key add -
sudo mkdir -p /etc/apt/sources.list.d/
echo 'echo "deb http://pkg.phz.fi/bionic ./" >> /etc/apt/sources.list.d/pkg.phz.fi.list' | sudo -s
sudo add-apt-repository http://pkg.phz.fi/
sudo apt-get update
sudo apt-get -y install phz-ric

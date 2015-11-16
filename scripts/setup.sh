#!/usr/bash

sudo apt-get update
sudo apt-get install -y libmagickwand-dev
sudo apt-get install -y golang-go

CONTAINED=$(grep "GOPATH" /home/vagrant/.profile | wc | gawk '//{ print $1; }')

if [[ $CONTAINED -eq "0" ]]; then
    echo "export PATH=/vagrant/bin:$PATH" >> /home/vagrant/.profile
    echo "export GOPATH=/vagrant" >> /home/vagrant/.profile
fi


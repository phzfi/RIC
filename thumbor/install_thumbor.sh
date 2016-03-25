#!/bin/bash

# Run this document with sudo

if [ -z $BASH_VERSION ] || [ "$EUID" -ne 0 ]; then
  sudo bash install_thumbor.sh
  exit 1
fi

apt-get install libssl-dev
apt-get install python-pip
apt-get install python-dev
apt-get install libcurl4-gnutls-dev
apt-get install libcurl4-openssl-dev
apt-get install libgraphicsmagick++1-dev libboost-python-dev
pip2 install graphicsmagick-engine
# Version is important to ensure graphicsmagick working"
pip2 install 'thumbor==5.2.1'

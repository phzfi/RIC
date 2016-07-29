#!/bin/bash

# Use with sudo

if [ -z $BASH_VERSION ] || [ "$EUID" -ne 0 ]; then
  sudo bash copy_test_images.sh
  exit 1
fi

mkdir /var/www
scp phzfi@ric.phz.fi:/var/www/*.jpg /var/www/


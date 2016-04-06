#!/bin/bash

# Runs with sudo. Make sure to run ./gradlew run once, before using the program.
if [ -z $BASH_VERSION ] || [ "$EUID" -ne 0 ]; then
  sudo bash install_cib.sh
  exit 1
fi

apt-get install openjdk-8-jdk
apt-get install Jmagick
update-ca-certificates -f
./gradlew run


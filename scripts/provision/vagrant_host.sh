#!/usr/bash

INSTALLED="$(which virtualbox | grep -c virtualbox)"

if [ ${INSTALLED} -eq 0 ]; then
    echo "Installing virtualbox"
    sudo apt-get install -y \
        virtualbox \
        virtualbox-dkms \
        virtualbox-guest-additions-iso
fi

INSTALLED="$(which vagrant | grep -c vagrant)"

if [ ${INSTALLED} -eq 0 ]; then
    echo "Installing vagrant"
    sudo apt-get install -y vagrant
fi


echo "Installing vagrant plugins"
for PLUGIN in "vagrant-vbguest" "vagrant-vbox-snapshot"; do
    vagrant plugin install "${PLUGIN}"
done



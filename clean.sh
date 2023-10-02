#!/bin/bash
#Clean test results etc

#Vagrant
~/workspace/vagrant/phz-vagrant-metadata/scripts/purge-boxes.sh

#Clean up sh2ju test results before next build to avoid cumulation
rm -fr results/*.xml

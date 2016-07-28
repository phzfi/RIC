#!/usr/bin/env bash

echo -e "\e[44mInstall global gulp for convenience...\e[0m"
sudo npm install -g gulp

echo -e "\e[44mSetting up 'riclib' npm package.json ...\e[0m"
cd /vagrant

# remove previously installed node modules
# to avoid issues between different node versions
rm -rf node_modules/

npm install --unsafe-perm

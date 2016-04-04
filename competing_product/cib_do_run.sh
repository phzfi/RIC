#!/bin/bash

screen -ls | awk -vFS='\t|[.]' '/CIB/ {system("screen -S "$2" -X quit")}'
cd ~/go/src/github.com/phzfi/RIC/competing_product/
screen -dmS CIB sh cib_loop.sh

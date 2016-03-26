#!/bin/bash

#shutdown RIC, CIB and Thumbor Instances
screen -ls | awk -vFS='\t|[.]' '/RIC/ {system("screen -S "$2" -X quit")}'
screen -ls | awk -vFS='\t|[.]' '/CIB/ {system("screen -S "$2" -X quit")}'
screen -ls | awk -vFS='\t|[.]' '/Thumbor/ {system("screen -S "$2" -X quit")}'

cd ~/go/src/github.com/phzfi/RIC/competing_product/

sudo sh clear_cache.sh

sh cib_do_run.sh

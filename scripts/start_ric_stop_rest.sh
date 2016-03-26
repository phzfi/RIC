#!/bin/bash

#shutdown RIC, CIB and Thumbor Instances
screen -ls | awk -vFS='\t|[.]' '/RIC/ {system("screen -S "$2" -X quit")}'
screen -ls | awk -vFS='\t|[.]' '/CIB/ {system("screen -S "$2" -X quit")}'
screen -ls | awk -vFS='\t|[.]' '/Thumbor/ {system("screen -S "$2" -X quit")}'

cd ~/go/src/github.com/phzfi/RIC/server/

echo "Clearing RIC temp files"
sudo rm -r /tmp/RICdiskcache

echo "Clearing cache"
sudo sh -c 'sync && echo 3 >/proc/sys/vm/drop_caches'

sh do_run.sh

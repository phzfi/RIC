#!/bin/bash

#shutdown RIC, CIB and Thumbor Instances
screen -ls | awk -vFS='\t|[.]' '/RIC/ {system("screen -S "$2" -X quit")}'
screen -ls | awk -vFS='\t|[.]' '/CIB/ {system("screen -S "$2" -X quit")}'
screen -ls | awk -vFS='\t|[.]' '/Thumbor/ {system("screen -S "$2" -X quit")}'

cd ~/go/src/github.com/phzfi/RIC/competing_product/

echo "Clearing CIB temp files"
sudo rm -r /tmp/tomcat*

echo "Clearing cache"
sudo sh -c 'sync && echo 3 >/proc/sys/vm/drop_caches'

sh cib_do_run.sh

#!/bin/bash

#shutdown RIC instances
screen -ls | awk -vFS='\t|[.]' '/RIC/ {system("screen -S "$2" -X quit")}'

cd ~/thumbor
thumbor --port=7777

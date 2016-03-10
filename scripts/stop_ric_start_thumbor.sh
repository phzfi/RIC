#!/bin/bash

#shutdown RIC instances
screen -ls | awk -vFS='\t|[.]' '/RIC/ {system("screen -S "$2" -X quit")}'

cd ~/thumbor

COUNT=`screen -ls | grep Thumbor | wc -l`
if [ $COUNT = "0" ]; then
  echo "Zero Thumbor screens found"
else
  echo $COUNT " Thumbor screens found"
  screen -ls | awk -vFS='\t|[.]' '/Thumbor/ {system("screen -S "$2" -X quit")}'
fi
screen -dmS Thumbor thumbor --port=7777

cd ~/go/src/github.com/phzfi/RIC/scripts/

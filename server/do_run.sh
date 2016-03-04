#!/bin/bash

COUNT=`screen -ls | grep RIC | wc -l`
if [ $COUNT = "0" ]; then
  echo "Zero RIC screens found"
else
  echo $COUNT " RIC screens found"
  screen -ls | awk -vFS='\t|[.]' '/RIC/ {system("screen -S "$2" -X quit")}'
fi
screen -dmS RIC ./server_loop.sh

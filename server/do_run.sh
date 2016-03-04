#!/bin/bash

TEST=`screen -ls | grep RIC | wc -l`
if [ $TEST = "0" ]; then
  echo "Zero RIC found"
else
  echo $TEST " found"
  screen -ls | awk -vFS='\t|[.]' '/RIC/ {system("screen -S "$2" -X quit")}'
fi
screen -dmS RIC ./server_loop.sh

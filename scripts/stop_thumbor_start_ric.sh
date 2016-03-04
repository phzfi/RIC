#!/bin/bash

#shutdown Thumbor instances
screen -ls | awk -vFS='\t|[.]' '/Thumbor/ {system("screen -S "$2" -X quit")}'

cd ~/go/src/github.com/phzfi/RIC/scripts/


COUNT=`screen -ls | grep RIC | wc -l`
if [ $COUNT = "0" ]; then
  echo "Zero RIC screens found"
else
  echo $COUNT " RIC screens found"
  screen -ls | awk -vFS='\t|[.]' '/RIC/ {system("screen -S "$2" -X quit")}'
fi
screen -dmS RIC ../server/do_run.sh

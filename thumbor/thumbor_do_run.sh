#!/bin/bash

screen -ls | awk -vFS='\t|[.]' '/Thumbor/ {system("screen -S "$2" -X quit")}'
cd  ~/go/src/github.com/phzfi/RIC/thumbor/
screen -dmS Thumbor sh thumbor_loop.sh

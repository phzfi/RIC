#!/bin/bash

screen -ls | awk -vFS='\t|[.]' '/RIC/ {system("screen -S "$2" -X quit")}'
screen -dmS RIC sh server_loop.sh

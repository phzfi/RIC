#!/bin/bash

if [ $# -ne 2 ] 
then
    echo 'Error: Wrong number of arguments. Usage: "./format_results.sh URLS_FILE OUTPUT_FILE"'
    exit
fi

RAW_FILE=$1
OUT_FILE=$2

python ./csv_formatter.py $RAW_FILE $OUT_FILE

#!/bin/bash
URLS_FILE=./urls.txt
RAW_FILE=./raw/$(date +%Y-%m-%d_%H-%M-%S).csv
OUT_FILE=./raw/$(date +%Y-%m-%d_%H-%M-%S).csv

./generate_report.sh $URLS_FILE $RAW_FILE
./format_results.sh $RAW_FILE $OUT_FILE

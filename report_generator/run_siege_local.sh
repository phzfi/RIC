#!/bin/bash

URLS_FILE=./urls_local.txt
RAW_FILE=./raw/ric_$(date +%Y-%m-%d_%H-%M-%S).txt
OUT_FILE=./results/ric_$(date +%Y-%m-%d_%H-%M-%S).csv
TMP=./temp/$(date +%Y-%m-%d_%H-%M-%S).tmp
SIEGE_CONF=./.siegerc
CONCURRENT=40
TIME="120s"

# Siege
siege -R $SIEGE_CONF --verbose --concurrent=$CONCURRENT -b --time=$TIME --log=$RAW_FILE --file=$URLS_FILE |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP
cat $TMP >> $RAW_FILE
rm $TMP

# Formatter
python csv_formatter.py $RAW_FILE $OUT_FILE

#!/bin/bash

URLS_FILE=./turls.txt
RAW_FILE=./raw/tumbor_test.txt
OUT_FILE=./results/tumbor_test.csv
TMP=./temp/text.tmp
SIEGE_CONF=./.siegerc
CONCURRENT=2
DELAY=3
TIME="40s"

# Siege

rm $RAW_FILE
siege -R $SIEGE_CONF --verbose --concurrent=$CONCURRENT --internet --delay=$DELAY --time=$TIME --log=$RAW_FILE --file=$URLS_FILE |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP
cat $TMP >> $RAW_FILE
rm $TMP

# Formatter
python csv_formatter.py $RAW_FILE $OUT_FILE

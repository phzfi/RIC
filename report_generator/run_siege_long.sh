#!/bin/bash

URLS_FILE=./urls.txt
OUT_FILE=./raw/$(date +%Y-%m-%d_%H-%M-%S).txt
TMP=./temp/temp$(date +%Y-%m-%d_%H-%M-%S).tmp
SIEGE_CONF=./.siegerc
CONCURRENT=2
DELAY=3
TIME="300s"

# Siege
siege -R $SIEGE_CONF --verbose --concurrent=$CONCURRENT --internet --delay=$DELAY --time=$TIME --log=$OUT_FILE --file=$URLS_FILE |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP 

cat $TMP >> $OUT_FILE
rm $TMP

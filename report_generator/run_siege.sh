#!/bin/bash
URLS_FILE=$1
OUT_FILE=$2

if [ $# -ne 2 ] 
then
    echo 'Error: Wrong number of arguments. Usage: "./generate_report.sh URLS_FILE OUTPUT_FILE"'
    exit
fi

TMP=./temp/temp$(date +%Y-%m-%d_%H-%M-%S).txt
LOG_FILE=./log/error.log
SIEGE_CONF=./.siegerc
CONCURRENT=2
DELAY=3
TIME="40s"

# Siege
siege -R $SIEGE_CONF --verbose --concurrent=$CONCURRENT --internet --delay=$DELAY --time=$TIME --log=$LOG_FILE --file=$URLS_FILE |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP 

cat $TMP >> $OUT_FILE
rm $TMP

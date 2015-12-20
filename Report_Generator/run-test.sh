#!/bin/bash
OUTNAME=./testResults/$(date +%Y-%m-%d_%H-%M-%S).csv
TMP=./temp/temp$(date +%Y-%m-%d_%H-%M-%S).txt

# Siege
siege -R .siegerc --verbose --concurrent=2 --internet --delay=3 --time=40s --log=$OUTNAME --file=urls.txt |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP 

cat $TMP >> $OUTNAME
rm $TMP

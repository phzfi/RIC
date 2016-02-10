#!/bin/bash

SEED=$1
CONCURRENT=$2
python urls_randomizer.py urls.txt $SEED
python urls_randomizer.py turls.txt $SEED


URLS_FILE=./urls.txt_temp.txt
TURLS_FILE=./turls.txt_temp.txt

# Siege settings
DELAY=2
TIME="480s"
SIEGE_CONF=./.siegerc

#RIC SIEGE
RAW_FILE=./raw/ric_$(date +%Y-%m-%d_%H-%M-%S).txt
RIC_OUT_FILE=./results/ric_$(date +%Y-%m-%d_%H-%M-%S).csv
TMP=./temp/$(date +%Y-%m-%d_%H-%M-%S).tmp


siege -R $SIEGE_CONF --verbose --concurrent=$CONCURRENT --delay=$DELAY --time=$TIME --log=$RAW_FILE --file=$URLS_FILE |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP
cat $TMP >> $RAW_FILE
rm $TMP

# Formatter
python csv_formatter.py $RAW_FILE $RIC_OUT_FILE

#TUMBOR SIEGE
RAW_FILE=./raw/tumbor_$(date +%Y-%m-%d_%H-%M-%S).txt
TUMBOR_OUT_FILE=./results/tumbor_$(date +%Y-%m-%d_%H-%M-%S).csv
TMP=./temp/$(date +%Y-%m-%d_%H-%M-%S).tmp


# Siege
siege -R $SIEGE_CONF --verbose --concurrent=$CONCURRENT --delay=$DELAY --time=$TIME --log=$RAW_FILE --file=$TURLS_FILE |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP
cat $TMP >> $RAW_FILE
rm $TMP

# Formatter
rm $URLS_FILE
rm $TURLS_FILE
python csv_formatter.py $RAW_FILE $TUMBOR_OUT_FILE
python csv_to_html.py $RIC_OUT_FILE $TUMBOR_OUT_FILE

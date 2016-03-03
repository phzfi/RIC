#!/bin/bash

if [ $# -lt 2 ]; then
    echo "Script usage: sh run_siege_tests.sh RANDOM_SEED CONCURRENT_USERS "
    exit 1
fi


SEED=$1
CONCURRENT=$2
python urls_randomizer.py urls_no_webp_local.txt $SEED
python urls_randomizer.py turls_no_webp_local.txt $SEED
python urls_randomizer.py jurls_no_webp_local.txt $SEED


URLS_FILE=./urls_no_webp_local.txt_temp.txt
TURLS_FILE=./turls_no_webp_local.txt_temp.txt
JURLS_FILE=./jurls_no_webp_local.txt_temp.txt

# Siege settings
DELAY=2
TIME="120s"
SIEGE_CONF=./.siegerc

#RIC SIEGE
RAW_FILE=./raw/ric_$(date +%Y-%m-%d_%H-%M-%S).txt
RIC_OUT_FILE=./results/ric_$(date +%Y-%m-%d_%H-%M-%S).csv
TMP=./temp/$(date +%Y-%m-%d_%H-%M-%S).tmp

# Siege
siege -R $SIEGE_CONF --verbose --concurrent=$CONCURRENT --delay=$DELAY --time=$TIME --log=$RAW_FILE --file=$URLS_FILE |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP
cat $TMP >> $RAW_FILE
rm $TMP

# Formatter
python csv_formatter.py $RAW_FILE $RIC_OUT_FILE

#TUMBOR SIEGE
RAW_FILE=./raw/thumbor_$(date +%Y-%m-%d_%H-%M-%S).txt
TUMBOR_OUT_FILE=./results/thumbor_$(date +%Y-%m-%d_%H-%M-%S).csv
TMP=./temp/$(date +%Y-%m-%d_%H-%M-%S).tmp


# Siege
siege -R $SIEGE_CONF --verbose --concurrent=$CONCURRENT --delay=$DELAY --time=$TIME --log=$RAW_FILE --file=$TURLS_FILE |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP
cat $TMP >> $RAW_FILE
rm $TMP

# Formatter
python csv_formatter.py $RAW_FILE $TUMBOR_OUT_FILE

# Competing Image Bank SIEGE
RAW_FILE=./raw/competitor_$(date +%Y-%m-%d_%H-%M-%S).txt
CIB_OUT_FILE=./results/competitor_$(date +%Y-%m-%d_%H-%M-%S).csv
TMP=./temp/$(date +%Y-%m-%d_%H-%M-%S).tmp

# Siege
siege -R $SIEGE_CONF --verbose --concurrent=$CONCURRENT --delay=$DELAY --time=$TIME --log=$RAW_FILE --file=$JURLS_FILE |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP
cat $TMP >> $RAW_FILE
rm $TMP

# Formatter
python csv_formatter.py $RAW_FILE $CIB_OUT_FILE

# Formatter
rm $URLS_FILE
rm $TURLS_FILE
rm $JURLS_FILE
python csv_to_html.py $RIC_OUT_FILE $TUMBOR_OUT_FILE $CIB_OUT_FILE


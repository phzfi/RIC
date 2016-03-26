#!/bin/bash

if [ $# -lt 3 ]; then
echo "Script usage: sh run_siege_same_request_amount.sh RANDOM_SEED CONCURRENT_USERS REQUESTS_PER_USER"
exit 1
fi


SEED=$1
CONCURRENT=$2
REQUESTS_PER_USER=$3
python urls_randomizer.py urls_no_webp.txt $SEED
python urls_randomizer.py turls_no_webp.txt $SEED


URLS_FILE=./urls_no_webp.txt_temp.txt
TURLS_FILE=./turls_no_webp.txt_temp.txt

# Siege settings
DELAY=1
SIEGE_CONF=./.siegerc

#RIC SIEGE
RAW_FILE=./raw/ric_$(date +%Y-%m-%d_%H-%M-%S).txt
RIC_OUT_FILE=./results/ric_$(date +%Y-%m-%d_%H-%M-%S).csv
TMP=./temp/$(date +%Y-%m-%d_%H-%M-%S).tmp

sh stop_thumbor_start_ric.sh
sleep 10
siege -R $SIEGE_CONF --verbose --concurrent=$CONCURRENT --delay=$DELAY -r$REQUESTS_PER_USER --log=$RAW_FILE --file=$URLS_FILE |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP
cat $TMP >> $RAW_FILE
rm $TMP

# Formatter
python csv_formatter.py $RAW_FILE $RIC_OUT_FILE

#TUMBOR SIEGE
RAW_FILE=./raw/tumbor_$(date +%Y-%m-%d_%H-%M-%S).txt
TUMBOR_OUT_FILE=./results/tumbor_$(date +%Y-%m-%d_%H-%M-%S).csv
TMP=./temp/$(date +%Y-%m-%d_%H-%M-%S).tmp
sh stop_ric_start_thumbor.sh
sleep 10s

# Siege
siege -R $SIEGE_CONF --verbose --concurrent=$CONCURRENT --delay=$DELAY -r$REQUESTS_PER_USER  --log=$RAW_FILE --file=$TURLS_FILE |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP
cat $TMP >> $RAW_FILE
rm $TMP


sh stop_thumbor_start_ric.sh


# Formatter
rm $URLS_FILE
rm $TURLS_FILE
python csv_formatter.py $RAW_FILE $TUMBOR_OUT_FILE
python csv_to_html.py constantRequestAmount.html $RIC_OUT_FILE $TUMBOR_OUT_FILE

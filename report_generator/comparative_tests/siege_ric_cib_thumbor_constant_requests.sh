#!/bin/bash


if [ $# -lt 3 ]; then
echo "Script usage: $0 RANDOM_SEED CONCURRENT_USERS REQUESTS_PER_USER"
exit 1
fi

# Checks if user started program properly, if not then starts it
if [ -z $BASH_VERSION ]; then
	exec bash "$0" "$@"
	exit $?
fi

# cd to report_generator
DIR="$( cd "$(dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"
cd ..

# Ric URLS randomiser
RURLS_FILE=siege_url_files/urls.txt
SEED=$1
python3 siege_url_files/urls_randomizer.py "$RURLS_FILE" $SEED
RURLS_FILE="${RURLS_FILE%.*}"_temp.txt

# CIB URLS randomiser
JURLS_FILE=siege_url_files/jurls.txt
SEED=$1
python3 siege_url_files/urls_randomizer.py "$JURLS_FILE" $SEED
JURLS_FILE="${JURLS_FILE%.*}"_temp.txt

# Thumbor URLS randomiser
TURLS_FILE=siege_url_files/turls.txt
SEED=$1
python3 siege_url_files/urls_randomizer.py "$TURLS_FILE" $SEED
TURLS_FILE="${TURLS_FILE%.*}"_temp.txt


# Siege settings
DELAY=1
SIEGE_CONF=./.siegerc
CONCURRENT=$2
REQUESTS_PER_USER=$3


#RIC SIEGE
RAW_FILE=./raw/ric_$(date +%Y-%m-%d_%H-%M-%S).txt
RIC_OUT_FILE=./results/ric_CR_"$SEED"_"$CONCURRENT"_"$REQUESTS_PER_USER"_$(date +%Y-%m-%d_%H-%M-%S).csv
TMP=./temp/$(date +%Y-%m-%d_%H-%M-%S).tmp


# Siege
siege -R $SIEGE_CONF --verbose --concurrent=$CONCURRENT --delay=$DELAY -r$REQUESTS_PER_USER --log=$RAW_FILE --file=$RURLS_FILE |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP
cat $TMP >> $RAW_FILE
rm $TMP


python3 csv_formatter.py $RAW_FILE $RIC_OUT_FILE


#CIB SIEGE
RAW_FILE=./raw/cib_$(date +%Y-%m-%d_%H-%M-%S).txt
CIB_OUT_FILE=./results/cib_CR_"$SEED"_"$CONCURRENT"_"$REQUESTS_PER_USER"_$(date +%Y-%m-%d_%H-%M-%S).csv
TMP=./temp/$(date +%Y-%m-%d_%H-%M-%S).tmp


# Siege
siege -R $SIEGE_CONF --verbose --concurrent=$CONCURRENT --delay=$DELAY -r$REQUESTS_PER_USER --log=$RAW_FILE --file=$JURLS_FILE |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP
cat $TMP >> $RAW_FILE
rm $TMP


python3 csv_formatter.py $RAW_FILE $CIB_OUT_FILE


#THUMBOR SIEGE
RAW_FILE=./raw/thumbor_$(date +%Y-%m-%d_%H-%M-%S).txt
THUMBOR_OUT_FILE=./results/thumbor_CR_"$SEED"_"$CONCURRENT"_"$REQUESTS_PER_USER"_$(date +%Y-%m-%d_%H-%M-%S).csv
TMP=./temp/$(date +%Y-%m-%d_%H-%M-%S).tmp


# Siege
siege -R $SIEGE_CONF --verbose --concurrent=$CONCURRENT --delay=$DELAY -r$REQUESTS_PER_USER --log=$RAW_FILE --file=$TURLS_FILE |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP
cat $TMP >> $RAW_FILE
rm $TMP

python3 csv_formatter.py $RAW_FILE $THUMBOR_OUT_FILE


# Formatter
rm $RURLS_FILE
rm $TURLS_FILE
rm $JURLS_FILE

python3 csv_to_html.py html_tables/ricCibThumborConstantRequestsResults.html $RIC_OUT_FILE $CIB_OUT_FILE $THUMBOR_OUT_FILE

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

# URLS randomiser
URLS_FILE=siege_url_files/turls.txt
SEED=$1
python3 siege_url_files/urls_randomizer.py "$URLS_FILE" $SEED
URLS_FILE="${URLS_FILE%.*}"_temp.txt

# Siege settings
DELAY=1
SIEGE_CONF=./.siegerc
CONCURRENT=$2
REQUESTS_PER_USER=$3


#THUMBOR SIEGE
RAW_FILE=./raw/thumbor_$(date +%Y-%m-%d_%H-%M-%S).txt
THUMBOR_OUT_FILE=./results/thumbor_CR_"$SEED"_"$CONCURRENT"_"$REQUESTS_PER_USER"_$(date +%Y-%m-%d_%H-%M-%S).csv
TMP=./temp/$(date +%Y-%m-%d_%H-%M-%S).tmp

# Siege
siege -R $SIEGE_CONF --verbose --concurrent=$CONCURRENT --delay=$DELAY -r$REQUESTS_PER_USER --log=$RAW_FILE --file=$URLS_FILE |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP
cat $TMP >> $RAW_FILE
rm $TMP

python3 csv_formatter.py $RAW_FILE $THUMBOR_OUT_FILE

# Formatter
rm $URLS_FILE

python3 csv_to_html.py html_tables/thumborConstantRequestsResults.html $THUMBOR_OUT_FILE

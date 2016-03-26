#!/bin/bash


if [ $# -lt 3 ]; then
echo "Script usage: $0 RANDOM_SEED CONCURRENT_USERS TIME_IN_SECONDS "
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
TIME=$3"s"


#THUMBOR SIEGE
RAW_FILE_BEFORE=./raw/thumbor-before_$(date +%Y-%m-%d_%H-%M-%S).txt
RAW_FILE_AFTER=./raw/thumbor-after_$(date +%Y-%m-%d_%H-%M-%S).txt
THUMBOR_OUT_FILE_BEFORE=./results/thumbor-before_CTA_"$SEED"_"$CONCURRENT"_"$TIME"_$(date +%Y-%m-%d_%H-%M-%S).csv
THUMBOR_OUT_FILE_AFTER=./results/thumbor-after_CTA_"$SEED"_"$CONCURRENT"_"$TIME"_$(date +%Y-%m-%d_%H-%M-%S).csv
TMP=./temp/$(date +%Y-%m-%d_%H-%M-%S).tmp

sh start_thumbor_stop_rest.sh

# Siege
siege -R $SIEGE_CONF --verbose --concurrent=$CONCURRENT --delay=$DELAY --time=$TIME --log=$RAW_FILE_BEFORE --file=$URLS_FILE |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP
cat $TMP >> $RAW_FILE_BEFORE
rm $TMP

TMP=./temp/$(date +%Y-%m-%d_%H-%M-%S).tmp

# Siege
siege -R $SIEGE_CONF --verbose --concurrent=$CONCURRENT --delay=$DELAY --time=$TIME --log=$RAW_FILE_AFTER --file=$URLS_FILE |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP
cat $TMP >> $RAW_FILE_AFTER
rm $TMP

python3 csv_formatter.py $RAW_FILE_BEFORE $THUMBOR_OUT_FILE_BEFORE
python3 csv_formatter.py $RAW_FILE_AFTER $THUMBOR_OUT_FILE_AFTER
# Formatter
rm $URLS_FILE

python3 csv_to_html.py html_tables/thumborConstantTimeResultsAuto.html $THUMBOR_OUT_FILE_BEFORE $THUMBOR_OUT_FILE_AFTER

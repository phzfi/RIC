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
URLS_FILE=/siege_url_files/urls_local.txt
SEED=$1
python /siege_url_files/urls_randomizer.py "$URLS_FILE" $SEED
URLS_FILE="${URLS_FILE%.*}"_temp.txt

# Siege settings
DELAY=1
SIEGE_CONF=./.siegerc
CONCURRENT=$2
TIME=$3"s"


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
python csv_to_html.py constantTimeResults.html $RIC_OUT_FILE $TUMBOR_OUT_FILE $CIB_OUT_FILE

>>>>>>> cc32638c629d0bec22b6ad21f4e0b48029be674f:report_generator/run_siege_tests_no_webp.sh

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
RURLS_FILE=siege_url_files/urls.txt
SEED=$1
python3 siege_url_files/urls_randomizer.py "$RURLS_FILE" $SEED
RURLS_FILE="${RURLS_FILE%.*}"_temp.txt

# Siege settings
DELAY=1
SIEGE_CONF=./.siegerc
CONCURRENT=$2
TIME=$3"s"


#RIC SIEGE
RAW_FILE=./raw/ric_$(date +%Y-%m-%d_%H-%M-%S).txt
RIC_OUT_FILE=./results/ric_CT_"$SEED"_"$CONCURRENT"_"$TIME"_$(date +%Y-%m-%d_%H-%M-%S).csv
TMP=./temp/$(date +%Y-%m-%d_%H-%M-%S).tmp

# Siege
siege -R $SIEGE_CONF --verbose --concurrent=$CONCURRENT --delay=$DELAY --time=$TIME --log=$RAW_FILE --file=$RURLS_FILE |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP
cat $TMP >> $RAW_FILE
rm $TMP

python3 csv_formatter.py $RAW_FILE $RIC_OUT_FILE

# Formatter
rm $RURLS_FILE

python3 csv_to_html.py html_tables/ricConstantTimeResults.html $RIC_OUT_FILE
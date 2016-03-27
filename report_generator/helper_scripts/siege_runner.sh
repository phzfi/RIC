#!/bin/bash

# This is a helper script intended to be called from siege_test.sh

URLS_FILE="$1"
SEED="$2"
REQUESTS_OR_TIME="$3"
CONCURRENT="$4"
LOCAL="$5"
AUTO="$6"
DELAY="$7"
SOFTWARE="$8"
RAW_FILE="$9"
OUT_FILE=${10}



# URL FILE RANDOMIZING
URLS_FILES_PATH=./siege_url_files/
URLS_FILE="$URLS_FILES_PATH""$URLS_FILE"

if [ "$LOCAL" == true ]; then
  URLS_FILE="${URLS_FILE%.*}"_local.txt
fi
python3 "$URLS_FILES_PATH"/urls_randomizer.py "$URLS_FILE" $SEED
URLS_FILE="${URLS_FILE%.*}"_temp.txt

# Siege settins/TMP
SIEGE_CONF=./.siegerc
TMP=./temp/$(date +%Y-%m-%d_%H-%M-%S).tmp

#If auto is true, run the auto scripts
if [ "$AUTO" == true ]; then
  if [ "$LOCAL" == true ]; then
    if [ "$SOFTWARE" == "cib" ]; then
      echo "Waiting 20s for cib to boot"
      sh ../scripts/start_cib_stop_rest.sh
      sleep 20s
    else
      echo "Waiting 10s for $SOFTWARE to boot"
      sh ../scripts/start_"$SOFTWARE"_stop_rest.sh
      sleep 10s
    fi
  else
    sh start_"$SOFTWARE"_stop_rest.sh
  fi
fi

echo "Running $SOFTWARE Siege!"

#Running Siege
siege -R $SIEGE_CONF --verbose $CONCURRENT $DELAY $REQUESTS_OR_TIME --log=$RAW_FILE --file=$URLS_FILE |
	 sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" > $TMP
cat $TMP >> $RAW_FILE
rm $TMP

python3 csv_formatter.py $RAW_FILE $OUT_FILE

rm $URLS_FILE

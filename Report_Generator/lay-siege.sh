#!/bin/bash
OUTNAME=./$(date +%Y-%m-%d_%H-%M-%S).log
TMP=$(mktemp)

# Siege
siege --verbose --concurrent=2 --internet --delay=3 --time=40s --log=$OUTNAME --file=urls.txt |
    sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]//g" |
    awk '/^HTTP/{ print $3, ",", $5, ",", $9, ",", $1, $2 }' |
    sed -r "s/ ?, ?/,/g" > $TMP

printf '\n\nRound trip time, size, path, response\n' >> $OUTNAME
cat $TMP >> $OUTNAME
rm -f $TMP

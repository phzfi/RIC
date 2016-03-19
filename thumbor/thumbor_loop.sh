#!/bin/bash

# A thumbor server loop that should be called from
# thumbor_do_run.sh found in the scripts folder.

# NOTE! This file should be in the folder ~/thumbor/
# with the thumbor.conf file

while true; do
    echo "Restaring thumbor"
    # 700 MB max
    ulimit -m 734003200
    thumbor --port=7777
    sleep 5
done

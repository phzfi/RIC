#!/bin/bash
ssh -t phzfi@ric.phz.fi 'sh ~/go/src/github.com/phzfi/RIC/scripts/start_ric_stop_rest.sh'
echo "Waiting 10s for ric to boot"
sleep 10s

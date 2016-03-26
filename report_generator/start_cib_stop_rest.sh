#!/bin/bash

ssh -t phzfi@ric.phz.fi 'sh ~/go/src/github.com/phzfi/RIC/scripts/start_cib_stop_rest.sh'
echo "Waiting 20s for cib to boot"
sleep 20s

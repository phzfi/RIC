#!/bin/bash

ssh -t phzfi@ric.phz.fi << EOF
  cd;
  cd go/src/github.com/phzfi/RIC/scripts;
  sh start_cib_stop_rest.sh;
EOF
echo "Waiting 20s for cib to boot"
sleep 20s

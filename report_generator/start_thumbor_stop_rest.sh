#!/bin/bash
ssh phzfi@ric.phz.fi << EOF
  cd;
  cd go/src/github.com/phzfi/RIC/scripts;
  sh start_thumbor_stop_rest.sh;
EOF
echo "Waiting 10s for thumbor to boot"
sleep 10s

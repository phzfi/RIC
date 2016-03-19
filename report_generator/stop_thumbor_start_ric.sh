#!/bin/bash
ssh phzfi@ric.phz.fi << EOF
  cd;
  cd go/src/github.com/phzfi/RIC/scripts;
  sh stop_thumbor_start_ric.sh;
EOF

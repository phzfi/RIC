#!/bin/bash
ssh phzfi@ric.phz.fi << EOF
	cd;
	cd go/src/github.com/phzfi/RIC/scripts;
	sh stop_ric_start_thumbor.sh;
EOF

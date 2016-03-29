#!/bin/bash
screen -ls | awk -vFS='\t|[.]' '/RIC/ {system("screen -S "$2" -X quit")}'
cd $GOPATH/src/github.com/phzfi/RIC/server
go build -tags profile && sh do_run.sh
echo "Waiting for ric to start"
sleep 4s
cd ../report_generator
sh siege_test.sh -l -d 0 -t 40 -c 200 & go tool pprof -svg http://localhost:6060/debug/pprof/profile > ../scripts/profiler_results/$(date +%Y-%m-%d_%H-%M-%S).svg

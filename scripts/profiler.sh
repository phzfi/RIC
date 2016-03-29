#!/bin/bash
screen -ls | awk -vFS='\t|[.]' '/RIC/ {system("screen -S "$2" -X quit")}'
sudo sh ~/go/src/github.com/phzfi/RIC/scripts/clear_cache.sh
cd ~/go/src/github.com/phzfi/RIC/server
echo "Building server"
go build -tags profile && sh do_run.sh
echo "Waiting for ric to start"
sleep 4s
cd ../report_generator
DATE=$(date +%Y-%m-%d_%H-%M-%S)
sh siege_test.sh -l -d 0 -t 35 -c 40 & go tool pprof -svg http://localhost:6060/debug/pprof/profile > ../scripts/profiler_results/fresh_$DATE.svg
sleep 30s
echo "Starting again"
sh siege_test.sh -l -d 0 -t 35 -c 40 & go tool pprof -svg http://localhost:6060/debug/pprof/profile > ../scripts/profiler_results/cached_$DATE.svg

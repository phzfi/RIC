#!/bin/bash

killall server
cd ../server/
sh do_run.sh


#create new screen with name foo screen -S foo
#to kill session foo  screen -X -S <sessionid> kill

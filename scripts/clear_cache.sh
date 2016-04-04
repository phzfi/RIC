#!/bin/bash

# To be run as sudo

echo "Clearing CIB temp files"
rm -r /tmp/tomcat*
echo "Clearing thumbor temp files"
rm -r /tmp/thumbor
echo "Clearing RIC temp files"
rm -r /tmp/RICdiskcache

echo "Clearing cache"
sh -c 'sync && echo 3 >/proc/sys/vm/drop_caches'

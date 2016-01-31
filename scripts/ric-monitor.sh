#!/bin/bash
TARGET="http://ric.phz.fi:8005/01.jpg?width=200&height=200"
NULL="/dev/null"

EXIT_CODE=$(curl \
    --connect-timeout 10 \
    --keepalive-time 30 \
    --max-time 60 \
    --output ${NULL} \
    --show-error \
    --write-out "%{http_code}" \
    ${TARGET} 2> ${NULL}
)

if [ ${EXIT_CODE} -eq "200" ]; then
    # Everything is OK!
    exit 0
elif [ ${EXIT_CODE} -eq "000" ]; then
    echo "Server is not responding"
    exit 1
else
    echo "Server reported error: HTTP_${EXIT_CODE}"
    exit 2
fi


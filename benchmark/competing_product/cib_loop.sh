#!/bin/bash

while true; do
    echo "Restaring cib"
    gradlew run
    sleep 5
done

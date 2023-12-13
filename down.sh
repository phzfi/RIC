#!/bin/bash
#Shutdown/destroy the dev env

function docker_down() {
    docker-compose down --remove-orphans
    exit 0
}

docker_down

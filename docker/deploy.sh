#!/bin/bash
#Deploy to prod/stg
ENV=$1

#TODO CHANGE ME
SERVICE_NAME=TODO_CHANGE_ME-$ENV


#No need to change anything below this line
export IMAGE_VERSION=$2
COMPOSE_FILE="docker-compose.${ENV}.yml"
#Password is in phz.kdbx
#DOCKER_REGISTRY_PASSWORD=${DOCKER_REGISTRY}
#export CONFIG_VERSION=`date +%Y%m%d%H%m`
export CONFIG_VERSION=$IMAGE_VERSION

if [ -z "$SERVICE_NAME" ] || [ -z "$IMAGE_VERSION" ]; then
    echo "Usage: ./deploy.sh <env> <version>, e.g. ./deploy.sh prod prod-124"
    exit 1
fi

if test -z "$DOCKER_REGISTRY_USERNAME"; then
    echo "ERROR: Building manually (outside Jenkins?). Please export DOCKER_REGISTRY_USERNAME to env from phz.kdbx"
    exit 1
fi
if test -z "$DOCKER_REGISTRY_PASSWORD"; then
    echo "ERROR: Building manually (outside Jenkins?). Please export DOCKER_REGISTRY_PASSWORD to env from phz.kdbx"
    exit 1
fi

#Do not deploy dev images
if test "$BUILD_ENV" == "stg" || test "$BUILD_ENV" == "prod"; then
    echo "Deploying $IMAGE_VERSION to $ENV"
    docker login docker-registry-in.phz.fi -u $DOCKER_REGISTRY_USERNAME -p $DOCKER_REGISTRY_PASSWORD

    export DOCKER_HOST=docker-swarm-master.in.phz.fi
    #docker stack rm $SERVICE_NAME
    docker stack deploy --with-registry-auth --compose-file docker-compose.$ENV.yml $SERVICE_NAME
    export DOCKER_HOST=
else
    echo "Skip deploy of $BUILD_ENV images"
fi


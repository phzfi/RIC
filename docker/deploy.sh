#!/bin/bash
#Deploy to prod/stg
BUILD_ENV=$1

SERVICE_NAME=ric-$ENV

#No need to change anything below this line
export IMAGE_VERSION=$2
COMPOSE_FILE="docker-compose.${ENV}.yml"
export CONFIG_VERSION=$IMAGE_VERSION

if [ -z "$SERVICE_NAME" ] || [ -z "$IMAGE_VERSION" ]; then
    echo "Usage: ./deploy.sh <env> <version>, e.g. ./deploy.sh prod prod-124"
    exit 1
fi

if ( [ "$BUILD_ENV" == "stg" ] || [ "$BUILD_ENV" == "prod" ] ) && test -z "$DOCKER_HUB_USERNAME"; then
    echo "ERROR: Building manually (outside Jenkins?). Please export DOCKER_HUB_USERNAME to env"
    exit 1
fi
if ( [ "$BUILD_ENV" == "stg" ] || [ "$BUILD_ENV" == "prod" ] ) && test -z "$DOCKER_HUB_PASSWORD"; then
    echo "ERROR: Building manually (outside Jenkins?). Please export DOCKER_HUB_PASSWORD to env"
    exit 1
fi

#Do not deploy dev images
if test "$BUILD_ENV" == "stg" || test "$BUILD_ENV" == "prod"; then
    echo "Deploying $IMAGE_VERSION to $ENV"
    docker login -u $DOCKER_HUB_USERNAME -p $DOCKER_HUB_PASSWORD

    #Deploy to Swarm
    export DOCKER_HOST=docker-swarm-master.in.phz.fi
    #docker stack rm $SERVICE_NAME
    docker stack deploy --with-registry-auth --compose-file docker-compose.$ENV.yml $SERVICE_NAME
    export DOCKER_HOST=
else
    echo "Skip deploy of $BUILD_ENV images"
fi


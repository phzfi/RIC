#!/bin/bash
#Note! Publish only staging and prod images, do not push dev images to docker-registry (but build dev version locally)

NAME=phz/ric

#No need to change anything below this line
BUILD_ENV=$1
VERSION=$2
TAG=$NAME:$VERSION

if ! ( [ "$BUILD_ENV" == "dev" ] || [ "$BUILD_ENV" == "stg" ] || [ "$BUILD_ENV" == "prod" ] ) || [ -z "$VERSION" ]; then
    echo "Usage: ./build.sh [dev|stg|prod] <version>, e.g. ./build.sh stg stg-123"
    exit 1
fi

if test -z $DOCKER_REGISTRY_USERNAME; then
    echo "Please provide docker-registry.in.phz.fi password in environment by variable DOCKER_REGISTRY_USERNAME from phz.kdbx, if you build this manually"
    exit 1
fi
if test -z $DOCKER_REGISTRY_PASSWORD; then
    echo "Please provide docker-registry.in.phz.fi password in environment by variable DOCKER_REGISTRY_PASSWORD from phz.kdbx, if you build this manually"
    exit 1
fi

export DOCKER_HOST=

echo "Building $TAG"

echo $BUILD_ENV-$VERSION > ../VERSION

docker login -u $DOCKER_HUB_USERNAME -p $DOCKER_HUB_PASSWORD \
    && docker build -f Dockerfile --build-arg BUILD_ENV="$BUILD_ENV" -t $TAG . \


# We don't want to push dev images to registry, just test if building works
if test "$BUILD_ENV" == "stg" || test "$BUILD_ENV" == "prod"; then
    docker tag $TAG $TAG \
    && docker push $TAG
else
    echo "Skip pushing $BUILD_ENV images to registry, since they are built locally"
fi

if test "$BUILD_ENV" == "prod"; then
    echo "Tagging and pushing $NAME:latest for prod"
    docker tag "$TAG" "$NAME":latest \
      && docker push $NAME:latest
fi


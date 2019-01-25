# Compile stage
FROM ubuntu:bionic AS ric-build-env

RUN apt-get update
RUN apt-get -y install software-properties-common golang-go git webp libwebp-dev file imagemagick libmagickwand-dev

# Set environment variables
ENV CGO_ENABLED 1
RUN echo 'export GOPATH=$HOME/go' >> $HOME/.bashrc
RUN echo 'export PATH=$PATH:$GOPATH/bin' >> $HOME/.bashrc

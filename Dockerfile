# Compile stage
FROM ubuntu:bionic AS ric-build-env

RUN apt-get update
RUN apt-get -y install software-properties-common
RUN apt-get update
RUN apt-get -y install golang-go

# Set environment variables
ENV GOPATH=$HOME
RUN echo 'export GOPATH=$HOME/go' >> $HOME/.bashrc
RUN echo 'export PATH=$PATH:$GOPATH/bin' >> $HOME/.bashrc

ENV CGO_ENABLED 1

# Update package handler packages
RUN apt-get update
RUN apt-get -y install git

# Install image generator tools
RUN apt-get -y install webp
RUN apt-get -y install libwebp-dev
RUN apt-get -y install file
RUN apt-get -y install imagemagick
RUN apt-get -y install libmagickwand-dev

#WORKDIR /root/go/src/github.com/phzfi/RIC/server

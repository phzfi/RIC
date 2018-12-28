# Compile stage
FROM ubuntu:bionic AS ric-build-env

RUN apt-get update
RUN apt-get -y install software-properties-common
#RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get update
RUN apt-get -y install golang-go

# Set environment variables
ENV GOPATH=$HOME
ENV PATH="${PATH}:$GOPATH/go/bin"
ENV CGO_ENABLED 1

# Update package handler packages
RUN apt-get update
RUN apt-get -y install git

# Install image preview generator tools
RUN apt-get -y install file
RUN apt-get -y install imagemagick
RUN apt-get -y install libmagickwand-dev

#ADD scripts/provision /provision
#
#RUN apt-get -y install wget
#
#RUN /provision/vagrant_setup_imagemagick.sh
#

WORKDIR /root/go/src/github.com/phzfi/RIC/server

# WIP
## Final stage
#FROM alpine:3.7 AS ric-prod-env
#
## For running binary files in alpine
#RUN apk add --no-cache libc6-compat
#
#EXPOSE 8005
#WORKDIR /
#COPY --from=ric-build-env /root/go/src/github.com/phzfi/RIC/server /server
#
#CMD ["/server/server"]

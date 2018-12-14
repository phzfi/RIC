# Compile stage
FROM golang:1.11.2 AS ric-build-env

# Set environment variables
ENV GOPATH=$HOME
ENV PATH="${PATH}:$GOPATH/go/bin"
ENV CGO_ENABLED 1

# Update package handler packages
RUN apt-get update

# Install image preview generator tools
RUN apt-get -y install file
RUN apt-get -y install imagemagick
RUN apt-get -y install libmagickwand-dev

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
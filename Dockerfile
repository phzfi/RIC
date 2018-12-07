# Compile stage
FROM golang:1.11.2 AS build-env

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

# GET go sources
WORKDIR /root/go/src/github.com/phzfi/RIC/server
COPY . .
CMD go get -t ./...
CMD go build
CMD ./server




#RUN go build -o RIC server
#
# Final stage
#FROM alpine:3.7
# For running binary files in alpine
#RUN apk add --no-cache libc6-compat
#EXPOSE 8005
##WORKDIR /
#COPY --from=build-env /server /
#CMD ["/server"]
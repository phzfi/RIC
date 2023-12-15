# Stage 1: Build the Go application
FROM golang:1.17 as go-builder

# Install ImageMagick dependencies
RUN apt-get update && apt-get install -y imagemagick libmagickwand-dev

# set go compiler options
ENV CGO_ENABLED=1
ENV GOOS=linux

WORKDIR /app

# copy local files to build container
COPY . .

# make directories
RUN mkdir -p /var/www
RUN mkdir -p /tmp

# initialise go project
RUN go mod init github.com/phzfi/RIC

# download necessary go libraries
RUN go get -t ./...
RUN go get -u ./...
RUN go mod download

# build Go application
RUN cd server; go build -v -tags debug -a -installsuffix cgo .


# Stage 2: Get certificates
FROM alpine:latest as certs
RUN apk update
RUN apk add --no-cache ca-certificates openssl-dev
RUN echo 'hosts: files dns' > /etc/nsswitch.conf


# Stage 3: Final stage, create the final image using Scratch
FROM scratch

# Copy SSL certs
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=certs /etc/nsswitch.conf /etc/nsswitch.conf

# Copy Go app
COPY --from=go-builder /app/server/server /ric-server
COPY --from=go-builder /app/server/config.ini /config.ini
COPY --from=go-builder /app/server/testimages/ /testimages/
COPY --from=go-builder /app/server/watermark.png /watermark.png
COPY --from=go-builder /app/server/testwm.png /testwm.png
COPY --from=go-builder /app/server/testresults /testresults
COPY --from=go-builder /app/server/config/testconfig.ini /config/testconfig.ini
COPY --from=go-builder /var/www /var/www
COPY --from=go-builder /tmp /tmp

# Copy any necessary libraries
COPY --from=go-builder /usr/lib/ /usr/lib/
COPY --from=go-builder /lib/x86_64-linux-gnu/ /lib/x86_64-linux-gnu/
COPY --from=go-builder /lib64/ /lib64/

ENTRYPOINT ["./ric-server"]

# Stage 1: Build the Go application
FROM golang:1.24-bookworm AS go-builder

# Install ImageMagick dependencies with OpenCL support
RUN apt-get update && apt-get install -y \
    imagemagick libmagickwand-dev \
    ocl-icd-libopencl1 opencl-headers clinfo \
    build-essential pkg-config libltdl-dev libjpeg-dev libpng-dev libtiff-dev libgif-dev libfreetype6-dev libwebp-dev libheif-dev libzip-dev \
    pocl-opencl-icd \
    fontconfig fonts-dejavu-core

# Build ImageMagick with OpenCL support
RUN cd /tmp && \
    wget -q https://imagemagick.org/archive/ImageMagick.tar.gz && \
    tar -xzf ImageMagick.tar.gz && \
    cd ImageMagick-* && \
    ./configure --prefix=/usr --enable-shared --enable-opencl --with-modules && \
    make -j$(nproc) && \
    make install && \
    ldconfig /usr/local/lib && \
    rm -rf /tmp/ImageMagick-*

ENV PATH="/usr/local/bin:$PATH"
ENV LD_LIBRARY_PATH="/usr/local/lib"

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
RUN cd server && go get -t ./...
RUN go mod vendor
RUN go mod download

# build Go application
RUN cd server; go build -v -tags debug -a -installsuffix cgo .


# Stage 2: Get certificates
FROM alpine:latest AS certs
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
COPY --from=go-builder /app/server/testimages/ /testimages/
COPY --from=go-builder /app/server/watermark.png /watermark.png
COPY --from=go-builder /app/server/testwm.png /testwm.png
COPY --from=go-builder /app/server/testresults /testresults
COPY --from=go-builder /var/www /var/www
COPY --from=go-builder /tmp /tmp

# Copy any necessary libraries
COPY --from=go-builder /usr/local/bin/ /usr/local/bin/
COPY --from=go-builder /usr/local/lib/ /usr/local/lib/
COPY --from=go-builder /usr/lib/ /usr/lib/
COPY --from=go-builder /lib/x86_64-linux-gnu/ /lib/x86_64-linux-gnu/
COPY --from=go-builder /lib64/ /lib64/
COPY --from=go-builder /etc/OpenCL/ /etc/OpenCL/
COPY --from=go-builder /etc/fonts /etc/fonts
COPY --from=go-builder /usr/share/fonts/truetype/dejavu /usr/share/fonts/truetype/dejavu
COPY --from=go-builder /var/cache/fontconfig /var/cache/fontconfig

ENV PATH="/usr/local/bin:/usr/bin:/bin"
ENV LD_LIBRARY_PATH="/usr/local/lib:/lib:/lib64:/usr/lib:/usr/lib/x86_64-linux-gnu"
ENV MAGICK_OCL_DEVICE=GPU

ENTRYPOINT ["./ric-server"]

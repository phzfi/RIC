# Use the latest Ubuntu image as the base
FROM ubuntu:latest

# Update the package repository
RUN apt-get update

# Install any necessary packages
# GLIB
RUN apt-get install -y libglib2.0-bin libglib2.0-dev
# DEVEL
RUN apt-get install -y automake autoconf gcc git g++ binutils make mercurial tar pkg-config vim wget bash git
# OPENCL
RUN apt-get install -y ocl-icd-opencl-dev opencl-headers ocl-icd-libopencl1  
# IMAGELIBS
RUN apt-get install -y libwebp-dev libtiff-dev libpng-dev libjpeg-dev liblqr-1-0-dev libltdl-dev
# GO
RUN apt-get install -y golang-go golang-golang-x-tools golang-golang-x-tools-dev \
  && rm -rf /var/lib/apt/lists/*



# Set the working directory in the container
WORKDIR /app

# Copy the local files to the container's working directory
COPY install_imagemagick.sh /app/

RUN chmod +x /app/install_imagemagick.sh

RUN /app/install_imagemagick.sh

COPY . .

RUN go mod init github.com/phzfi/RIC
RUN go get -t ./...
RUN go list -u -m all
RUN go get -u ./...

# build with debug
# RUN cd server; go build -tags debug .

# build production
RUN cd server; go build .

RUN mkdir -p /var/www

WORKDIR /app/server

# Command to run when the container starts
CMD ["./server"]

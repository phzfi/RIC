#!/bin/bash

IMAGELIBS=libwebp-dev libtiff5-dev libpng12-dev libjpeg-dev liblqr-1-0-dev
OPENCL=ocl-icd-opencl-dev opencl-headers ocl-icd-libopencl1

# Install the absolute requirements
apt-get update
apt-get install -y ${IMAGELIBS}
apt-get install -y ${OPENCL}
apt-get install -y git gcc binutils automake autoconf

source "vagrant_setup_imagemagick.sh"
source "vagrant_setup_go.sh"


#!/bin/bash

GLIB="libglib2.0-bin libglib2.0-dev"
IMAGELIBS="libwebp-dev libtiff5-dev libpng12-dev libjpeg-dev liblqr-1-0-dev"
OPENCL="ocl-icd-opencl-dev opencl-headers ocl-icd-libopencl1"
DEVEL="automake autoconf gcc git g++ binutils make mercurial tar pkg-config vim wget"

# Install the absolute requirements
apt-get update
apt-get install -y ${GLIB}
apt-get install -y ${DEVEL}
apt-get install -y ${OPENCL}
apt-get install -y ${IMAGELIBS}

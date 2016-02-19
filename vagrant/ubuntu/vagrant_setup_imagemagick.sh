#!/bin/bash

# ImageMagick 6.9.3-4 (2016-02-12)
REPO="http://git.imagemagick.org/repos/ImageMagick.git"
RELEASE="d21fb1eaf6e444ecd6228f2a58d6d0e24692f53f"
SOURCE="imagemagick-source"
BUILD="imagemagick-build"

cd /tmp
git clone "${REPO}" "${SOURCE}"
if [ ! -d "${SOURCE}" ]; then
	echo "Failed to clone ImageMagick repository!"
	exit 1
fi

mkdir -p "${BUILD}"
cd imagemagick-source
git checkout "${RELEASE}"
cd "../${BUILD}"

../${SOURCE}/configure \
	--prefix=/usr \
	--enable-opencl

make
make install


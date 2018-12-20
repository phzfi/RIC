#!/bin/bash

ORIGIN="https://www.imagemagick.org/download/ImageMagick-6.9.10-16.tar.gz"
#ORIGIN="https://www.imagemagick.org/download/ImageMagick.tar.gz"
FILE="imagemagick.tar.gz"
SOURCE="imagemagick-source"
BUILD="imagemagick-build"

cd /tmp

mkdir -p "${SOURCE}"
mkdir -p "${BUILD}"

wget --no-verbose --tries=10 "${ORIGIN}" -O "${FILE}"
if [ $? -ne 0 ]; then
	echo "Downloading imagemagick failed!"
	exit 1
fi

tar --directory="${SOURCE}" --strip-components=1 -xf "${FILE}"
if [ $? -ne 0 ]; then
	echo "Extracting imagemagick failed!"
	exit 1
fi


cd "${BUILD}"
../${SOURCE}/configure \
	--prefix=/usr \
	--enable-opencl

make -j2
if [ $? -ne 0 ]; then
	echo "Building imagemagick failed!"
	exit 1
fi

make install
if [ $? -ne 0 ]; then
	echo "Installing imagemagick failed!"
	exit 1
fi


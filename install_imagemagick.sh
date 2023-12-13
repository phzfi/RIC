#!/bin/bash
set -e

# ImageMagick 6.9.13-1 (2023-12-10)
ORIGIN="https://imagemagick.org/archive/ImageMagick-6.9.13-1.tar.bz2"
# RELEASE="d21fb1eaf6e444ecd6228f2a58d6d0e24692f53f"
# ORIGIN="http://git.imagemagick.org/repos/ImageMagick/repository/archive.tar.bz2?ref=${RELEASE}"
FILE="imagemagick.tar.bz2"
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

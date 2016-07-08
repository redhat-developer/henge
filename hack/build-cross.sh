#!/bin/bash

#
# Script for cross compiling Henge for multiple platforms.
#

HENGE_ROOT=$(dirname "${BASH_SOURCE}")/..
BUILD_DIR=${HENGE_ROOT}/bin

# Platforms and architectures we want to build for
PLATFORMS=${1:-"darwin linux windows"}
ARCHS=${2:-"amd64"}

echo "Staring builds"

for platform in $PLATFORMS; do
  for arch in $ARCHS; do
    echo " Building for ${platform}/${arch}"
    OUTPUT_FILE=${BUILD_DIR}/${platform}/${arch}/henge
    GOOS=$platform GARCH=$arch go build -o ${OUTPUT_FILE}
  done
done

echo "Builds finished"
echo "Binaries are in `readlink -f ${BUILD_DIR}`/"

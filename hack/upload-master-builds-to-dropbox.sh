#!/bin/bash

# This uploads cross compiled binaries to Dropbox
# https://www.dropbox.com/sh/lcz3o5o0i0gv6fw/AADMZrz7URk0R6qkKnyE7Hs7a?dl=0


# This should be set outside of this script in some secure way
# DROPBOX_TOKEN=

SHORT_COMMIT=`echo $TRAVIS_COMMIT | cut -c1-8`

PLATFORMS="darwin linux windows"
ARCHS="amd64"


upload_file() {
  SOURCE=$1
  TARGET=$2

  curl -X POST https://content.dropboxapi.com/2/files/upload \
    --header "Authorization: Bearer ${DROPBOX_TOKEN}" \
    --header "Content-Type: application/octet-stream" \
    --header "Dropbox-API-Arg: {\"path\":\"${TARGET}\",\"autorename\":true}" \
    --data-binary @"${SOURCE}"

}

#record commit id
echo $TRAVIS_COMMIT > "commitid"


for platform in $PLATFORMS; do
  for arch in $ARCHS; do
    # upload binary
    upload_file "bin/${platform}/${arch}/henge" "/HengeMasterBuilds/${platform}/${arch}/henge-latest"
    upload_file "build-from" "/HengeMasterBuilds/${platform}/${arch}/commitid"
  done
done

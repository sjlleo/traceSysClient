#!/usr/bin/env bash

set -e

DIST_PREFIX="traceClient"
DEBUG_MODE=${2}
TARGET_DIR="dist"
PLATFORMS="darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 linux/s390x linux/mips windows/amd64 windows/arm64"

rm -rf ${TARGET_DIR}
mkdir ${TARGET_DIR}

for pl in ${PLATFORMS}; do
    export GOOS=$(echo ${pl} | cut -d'/' -f1)
    export GOARCH=$(echo ${pl} | cut -d'/' -f2)
    export TARGET=${TARGET_DIR}/${DIST_PREFIX}_${GOOS}_${GOARCH}
    if [ "${GOOS}" == "windows" ]; then
        export TARGET=${TARGET_DIR}/${DIST_PREFIX}_${GOOS}_${GOARCH}.exe
    fi

    echo "build => ${TARGET}"
    if [ "${DEBUG_MODE}" == "debug" ]; then
        go build -trimpath -gcflags "all=-N -l" -o ${TARGET} \
            -ldflags    "-X 'main.version=${BUILD_VERSION}' \
                        -X 'main.buildDate=${BUILD_DATE}' \
                        -X 'main.commitID=${COMMIT_SHA1}'\
                        -w -s"
    else
        go build -trimpath -o ${TARGET} \
            -ldflags    "-X 'main.version=${BUILD_VERSION}' \
                        -X 'main.buildDate=${BUILD_DATE}' \
                        -X 'main.commitID=${COMMIT_SHA1}'\
                        -w -s"
    fi
done

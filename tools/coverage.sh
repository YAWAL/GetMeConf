#!/bin/bash

COVERAGE_DIR="${GOPATH}/src/github.com/YAWAL/GetMeConf/tools/coverage"

PKG_LIST=$(go list ./... | grep -v /vendor/) ;

mkdir -p "$COVERAGE_DIR";

for PACKAGE in ${PKG_LIST}; do
    go test -covermode=count -coverprofile "${COVERAGE_DIR}/${PACKAGE##*/}.cov" "$PACKAGE" ;
done ;


if [ "$1" == "html" ]; then
    mkdir -p "$COVERAGE_DIR/html";
    for FILE in $COVERAGE_DIR/*.cov; do
        FILE_NAME=$(basename "$FILE" .cov)
        go tool cover -html="${FILE}" -o "${COVERAGE_DIR}/html/${FILE_NAME}.html" ;
    done
fi
#!/bin/bash
COVERAGE_DIR=${GOPATH}/src/github.com/YAWAL/GetMeConf/tools
COVERAGE_DIR="${COVERAGE_DIR:-coverage}"
PKG_LIST=${go list ./... | grep -v /vendor/)

mkdir -p "$COVERAGE_DIR";

for package in ${PKG_LIST}; do
    go test -covermode=count -coverprofile "${COVERAGE_DIR}/${package##*/}.cov" "$package" ;
done ;

echo 'mode: count' > "${COVERAGE_DIR}"/coverage.cov ;

if [ "$1" == "html" ]; then
    go tool cover -html="${COVERAGE_DIR}"/coverage.cov -o coverage.html ;
fi
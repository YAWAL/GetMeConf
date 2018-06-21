# force to use bash
SHELL = /bin/bash

include .env
export

RICHGO := $(shell PATH=$(GOPATH)/bin:$(PATH) command -v richgo 2> /dev/null)
ifdef RICHGO
GO := $(RICHGO)
else
$(shell GOPATH=$(GOPATH) go get -u github.com/kyoh86/richgo)
GO := $(GOPATH)/bin/richgo
endif

all: dependencies build

.PHONY: build run build dependencies dep tests race

build:
	echo "Build"
	go build -o ${GOPATH}/src/github.com/YAWAL/GetMeConf/bin/service .

run:
	echo "Running service"
	go build -o ${GOPATH}/src/github.com/YAWAL/GetMeConf/bin/service .
	./bin/service

dependencies:
	echo "Installing dependencies"
	dep ensure

install dep:
	echo    "Installing dep"
	curl    https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

tests:
	echo "Tests"
	go test ./service
	go test ./repository

race:
	echo "Race tests"
	go test ./service -race
	go test ./repository -race

lint:
	echo "Running linters"
	gometalinter --vendor --tests --skip=mock --exclude='_gen.go' --deadline=1500s --checkstyle --sort=linter ./... > static-analysis.xml

docker-build:
	CC=$(which musl-gcc) go build --ldflags '-w -linkmode external -extldflags "-static"' -o ${GOPATH}/src/github.com/YAWAL/GetMeConf/bin/service . && \
	docker build -t configservice . && \
	docker run --net=${DOCKER_NET_DRIVER} -p ${SERVICE_PORT}:${SERVICE_PORT} --env-file .env configservice

clean:
	echo "Removing previous build"
	rm -rf ${GOPATH}/src/github.com/YAWAL/GetMeConf/bin/service

coverage:
	./tools/coverage.sh;

coveragehtml:
	./tools/coverage.sh html;
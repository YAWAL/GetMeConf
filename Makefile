include .env
export

all: dependencies build

.PHONY: build
build:
	echo "Build"
	go build -o ${GOPATH}/src/github.com/YAWAL/GetMeConf/bin/service ./service

.PHONY: run
run:
	echo "Running service"
	go build -o ${GOPATH}/src/github.com/YAWAL/GetMeConf/bin/service ./service
	./bin/service

.PHONY: dependencies
dependencies:
	echo "Installing dependencies"
	dep ensure

install dep:
	echo    "Installing dep"
	curl    https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

.PHONY: tests
tests:
	echo "Tests"
	go test ./service
	go test ./repository

race:
	echo "Race tests"
	go test ./service -race
	go test ./repository -race

docker-build:
	docker run --net=${DOCKER_NET_DRIVER} -d consul && \
	CC=$(which musl-gcc) go build --ldflags '-w -linkmode external -extldflags "-static"' -o ${GOPATH}/src/github.com/YAWAL/GetMeConf/bin/service ./service && \
	docker build -t configservice . && \
	docker run --net=${DOCKER_NET_DRIVER} -p ${SERVICE_PORT}:${SERVICE_PORT} --env-file .env configservice

clean:
	echo "Removing previous build"
	rm -rf ${GOPATH}/src/github.com/YAWAL/GetMeConf/bin/service

coverage:
	echo "Test coverage"
	./tools/coverage.sh;

coveragehtml:
	echo "Test coverage"
	./tools/coverage.sh html;
include .env
export

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
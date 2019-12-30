BUILD=`date +%FT%T%z`
COMMIT=`git rev-parse HEAD`
VER=1.1.3

LDFLAGS=-ldflags " -s -X main.AppVersion=${VER} -X main.BuildDate=${BUILD} -X main.GitCommit=${COMMIT}"

build:
	rm -rf dist
	mkdir dist
	go build ${LDFLAGS} -o ./dist/pubsub .

build-race:
	rm -rf dist
	mkdir dist
	go build -race ${LDFLAGS} -o ./dist/pubsub .

build-all:
	rm -rf dist
	mkdir dist
	go build -race ${LDFLAGS} -o ./dist/pubsub .
	go build -race -o ./dist/pub ./bench/pub/*.go
	go build -race -o ./dist/sub ./bench/sub/*.go

run:
	go run --race main.go -c ./config.toml


start: build-race
	chmod -R +x ./dist/pubsub
	GODEBUG=gctrace=1 ./dist/pubsub

clean:
	rm -rf dist
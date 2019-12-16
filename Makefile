BUILD=`date +%FT%T%z`
COMMIT=`git rev-parse HEAD`
VER=1.1.3

LDFLAGS=-ldflags " -s -X main.AppVersion=${VER} -X main.BuildDate=${BUILD} -X main.GitCommit=${COMMIT}"

build:
	rm -rf dist
	mkdir dist
	go build ${LDFLAGS} -o ./dist/pubsub .

race-build:
	rm -rf dist
	mkdir dist
	go build -race ${LDFLAGS} -o ./dist/pubsub .

run:
	go run --race main.go start -c ./config.toml


start: race-build
	chmod -R +x ./dist/pubsub
	GODEBUG=gctrace=1 ./dist/pubsub start

clean:
	rm -rf dist
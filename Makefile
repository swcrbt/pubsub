BUILD=`date +%FT%T%z`
COMMIT=`git rev-parse HEAD`
VER=2.0.1

LDFLAGS=-ldflags " -s -X main.AppVersion=${VER} -X main.BuildDate=${BUILD} -X main.GitCommit=${COMMIT}"

build:
	rm -rf dist
	mkdir dist
	go build ${LDFLAGS} -o ./dist/issued-service .

run:
	go run --race main.go -c ./config.toml
	
clean:
	rm -rf dist
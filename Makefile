BUILD=`date +%FT%T%z`
COMMIT=`git rev-parse HEAD`
VER=1.0.1

LDFLAGS=-ldflags " -s -X main.AppVersion=${VER} -X main.BuildDate=${BUILD} -X main.GitCommit=${COMMIT}"

build:
	rm -rf dist
	mkdir dist
	go build ${LDFLAGS} -o ./dist/issue .

run:
	go run --race main.go start -c ./config.toml

start: build
	chmod -R +x ./dist/issue
	./dist/issue start

clean:
	rm -rf dist
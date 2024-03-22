# Merlin version
VERSION=$(shell cat version/version.go |grep "const Version ="|cut -d"\"" -f2)
BUILD=$(shell git rev-parse HEAD)
DIR=bin/v${VERSION}/${BUILD}
DEST=.

# Go build flags
LDFLAGS=-ldflags '-X github.com/Ne0nd0g/merlin-cli/version.Build=${BUILD}'

default:
	go build ${LDFLAGS} -o merlinCLI main.go

# Build all
all: darwin linux windows

# Compile Server - Darwin x64
darwin:
	export GOOS=darwin;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/merlinCLI-Darwin-x64 main.go

# Compile Server - Linux x64
linux:
	export GOOS=linux;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/merlinCLI-Linux-x64 main.go

# Compile Server - Windows x64
windows:
	export GOOS=windows;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/merlinCLI-Windows-x64.exe main.go

move:
	mv ${DIR}/merlinCLI-* ${DEST}

release: all move

clean:
	rm -rf merlinCLI-*
	rm -rf ${DIR}/merlinCLI-*
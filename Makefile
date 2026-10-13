# Merlin version
VERSION=$(shell cat version/version.go |grep "const Version ="|cut -d"\"" -f2)
BUILD=$(shell git rev-parse HEAD)
DIR=bin/v${VERSION}/${BUILD}
DEST=.

# Go build flags
LDFLAGS=-ldflags '-X github.com/Ne0nd0g/merlin-cli/version.Build=${BUILD}'

# Build fully static, pure-Go binaries for every target. merlin-cli uses no cgo,
# so this only forces static linking + Go's native net/user resolvers and makes
# builds deterministic regardless of which C toolchains exist on the build host.
export CGO_ENABLED := 0

default:
	go build ${LDFLAGS} -o merlinCLI main.go

# Build all
all: darwin darwin-arm64 linux windows arm arm64

# Compile Server - Darwin x64
darwin:
	export GOOS=darwin;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/merlinCLI-Darwin-x64 main.go

# Compile Darwin arm64 (Apple Silicon)
darwin-arm64:
	export GOOS=darwin;export GOARCH=arm64;go build ${LDFLAGS} -o ${DIR}/merlinCLI-Darwin-arm64 main.go

# Compile Server - Linux x64
linux:
	export GOOS=linux;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/merlinCLI-Linux-x64 main.go

# Compile Server - Windows x64
windows:
	export GOOS=windows;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/merlinCLI-Windows-x64.exe main.go

# Compile Server - Linux arm (armv7)
arm:
	export GOOS=linux;export GOARCH=arm;export GOARM=7;go build ${LDFLAGS} -o ${DIR}/merlinCLI-Linux-arm main.go

# Compile Server - Linux arm64
arm64:
	export GOOS=linux;export GOARCH=arm64;go build ${LDFLAGS} -o ${DIR}/merlinCLI-Linux-arm64 main.go

move:
	mv ${DIR}/merlinCLI-* ${DEST}

release: all move

clean:
	rm -rf merlinCLI-*
	rm -rf ${DIR}/merlinCLI-*
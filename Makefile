gotag=1.20.14-bullseye	# support win7
gotag_modern=1.23.0-bullseye

commit=$(shell git rev-parse HEAD)

dockercmd=docker run --rm -v $(CURDIR):/usr/src/qbtchangetracker -w /usr/src/qbtchangetracker
buildtags = -tags forceposix
buildenvs = -e CGO_ENABLED=0
version = 1.999
ldflags = -ldflags="-X 'main.version=$(version)' -X 'main.commit=$(commit)' -X 'main.buildImage=golang:$(gotag)'"
ldflags_modern = -ldflags="-X 'main.version=$(version)' -X 'main.commit=$(commit)' -X 'main.buildImage=golang:$(gotag_modern)'"

.PHONY: all tests build

all: | tests build

tests:
	$(dockercmd) golang:$(gotag) go test $(buildtags) ./...

build: windows linux darwin

windows:
	$(dockercmd) $(buildenvs) -e GOOS=windows -e GOARCH=amd64 golang:$(gotag) go build -v $(buildtags) $(ldflags) -o qbtchangetracker_v$(version)_amd64.exe
	$(dockercmd) $(buildenvs) -e GOOS=windows -e GOARCH=arm64 golang:$(gotag_modern) go build -v $(buildtags) $(ldflags_modern) -o qbtchangetracker_v$(version)_arm64.exe
	$(dockercmd) $(buildenvs) -e GOOS=windows -e GOARCH=386 golang:$(gotag) go build -v $(buildtags) $(ldflags) -o qbtchangetracker_v$(version)_i386.exe

linux:
	$(dockercmd) $(buildenvs) -e GOOS=linux -e GOARCH=amd64 golang:$(gotag_modern) go build -v $(buildtags) $(ldflags_modern) -o qbtchangetracker_v$(version)_amd64_linux
	$(dockercmd) $(buildenvs) -e GOOS=linux -e GOARCH=386 golang:$(gotag_modern) go build -v $(buildtags) $(ldflags_modern) -o qbtchangetracker_v$(version)_i386_linux

darwin:
	$(dockercmd) $(buildenvs) -e GOOS=darwin -e GOARCH=amd64 golang:$(gotag_modern) go build -v $(buildtags) $(ldflags_modern) -o qbtchangetracker_v$(version)_amd64_macos
	$(dockercmd) $(buildenvs) -e GOOS=darwin -e GOARCH=arm64 golang:$(gotag_modern) go build -v $(buildtags) $(ldflags_modern) -o qbtchangetracker_v$(version)_arm64_macos

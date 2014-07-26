.PHONY: all binary library libgospace gospace deps test format

GOPATH=$(PWD)
GOBIN=$(PWD)/bin

# ensure build artifacts end up in the current directory
export GOPATH
export GOBIN

default: binary

all: binary

binary: deps gospace

library: deps libgospace

libgospace:
	go install gospace

gospace: library
	go install cli/gospace

deps:
	go get -d -v ./...

test: deps
	go test ./...

format:
	go fmt gospace
	go fmt cli/gospace
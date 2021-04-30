NAME = alm-mqtt-module
BIN_DIR ?= ../bin
VERSION ?= $(shell git describe --match=NeVeRmAtCh --always --abbrev=40 --dirty)
GO_LDFLAGS = -tags 'netgo osusergo static_build' -ldflags "-X alm-mqtt-module/internal/version.Version=$(VERSION)"

all: test build

build:
	GOOS=linux go build $(GO_LDFLAGS) -o ${BIN_DIR}/${NAME} main.go

test:
	go test ./...

example:
	cd example/client && make

clean:
	rm -f ${BIN_DIR}/${NAME}
	cd example/client && make clean

.PHONY: clean example

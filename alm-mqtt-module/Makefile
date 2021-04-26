NAME = alm-mqtt-module
BIN_DIR ?= ../bin
VERSION ?= $(shell git describe --match=NeVeRmAtCh --always --abbrev=40 --dirty)
GO_LDFLAGS = -ldflags "-X alm-mqtt-module/internal/version.Version=$(VERSION)"

all: build

build:
	GOOS=linux go build $(GO_LDFLAGS) -o ${BIN_DIR}/${NAME} main.go

test:
	go test ./...

example:
	cd example && make

clean:
	rm -f ${BIN_DIR}/${NAME}
	cd example && make clean

.PHONY: clean example

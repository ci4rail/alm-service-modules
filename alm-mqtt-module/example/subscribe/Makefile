NAME = alm-mqtt-module-subscribe
BIN_DIR ?= ../../../bin
GO_LDFLAGS = -tags 'netgo osusergo static_build'

all: test build

build:
	GOOS=linux go build $(GO_LDFLAGS) -o ${BIN_DIR}/${NAME} subscribe.go

test:
	go test ./...

clean:
	rm -f ${BIN_DIR}/${NAME}

.PHONY: clean

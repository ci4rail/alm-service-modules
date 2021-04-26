NAME = alm-mqtt-module-example
BIN_DIR ?= ../../bin
GO_LDFLAGS = -tags 'netgo osusergo static_build'

all: build

build:
	GOOS=linux go build $(GO_LDFLAGS) -o ${BIN_DIR}/${NAME} example.go

test:
	go test ./...

clean:
	rm -f ${BIN_DIR}/${NAME}

.PHONY: clean

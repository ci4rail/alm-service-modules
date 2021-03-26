FROM golang:1.16 AS build
WORKDIR /go/src/
COPY . /go/src/alm-location-module
ENV CGO_ENABLED=0
ENV GOPATH=/go
ARG VERSION=dev

WORKDIR /go/src/alm-location-module
RUN VERSION=${VERSION} BIN_DIR=/install make

FROM alpine:3.12
COPY --from=build /install/alm-location-module /alm-location-module
ENTRYPOINT ["/alm-location-module"]

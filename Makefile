SHORT_NAME ?= deisrel

VERSION ?= git-$(shell git rev-parse --short HEAD)
LDFLAGS := "-s -X main.version=${VERSION}"

HOST_OS := $(shell uname)
ifeq ($(HOST_OS),Darwin)
	GOOS=darwin
else
	GOOS=linux
endif

REPO_PATH := github.com/deis/${SHORT_NAME}
DEV_ENV_IMAGE := quay.io/deis/go-dev:0.13.0
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}
DEV_ENV_PREFIX := docker run --rm -e GO15VENDOREXPERIMENT=1 -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_CMD := ${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE}

DEIS_BINARY_NAME ?= ./deis

bootstrap:
	${DEV_ENV_CMD} glide install

build:
	${DEV_ENV_PREFIX} -e GOOS=${GOOS} ${DEV_ENV_IMAGE} go build -a -installsuffix cgo -ldflags ${LDFLAGS} -o deisrel .

test:
	${DEV_ENV_CMD} sh -c 'go test -v $$(glide nv)'

build-cli-cross:
	${DEV_ENV_CMD} gox -output="bin/${SHORT_NAME}-${VERSION}-{{.OS}}-{{.Arch}}"
	${DEV_ENV_CMD} gox -output="bin/${SHORT_NAME}-latest-{{.OS}}-{{.Arch}}"

dist: build-cli-cross

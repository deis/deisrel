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
DEV_ENV_IMAGE := quay.io/deis/go-dev:0.14.0
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}
DEV_ENV_PREFIX := docker run --rm -e GO15VENDOREXPERIMENT=1 -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_CMD := ${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE}

DEIS_BINARY_NAME ?= ./deis

GO_BUILD_CMD := go build -a -installsuffix cgo -ldflags ${LDFLAGS} -o deisrel .
GO_TEST_CMD := go test -v $$(glide nv)

CROSS_COMPILE_OPTS := -os="linux darwin " -arch="amd64"

bootstrap:
	${DEV_ENV_CMD} glide install

build:
	${GO_BUILD_CMD}

build-docker:
	${DEV_ENV_PREFIX} -e GOOS=${GOOS} ${DEV_ENV_IMAGE} ${GO_BUILD_CMD}

test:
	${GO_TEST_CMD}

test-docker:
	${DEV_ENV_CMD} sh -c '${GO_TEST_CMD}'

build-cli-cross:
	${DEV_ENV_CMD} gox ${CROSS_COMPILE_OPTS} -output="bin/${VERSION}/${SHORT_NAME}-${VERSION}-{{.OS}}-{{.Arch}}"
	${DEV_ENV_CMD} gox ${CROSS_COMPILE_OPTS} -output="bin/${SHORT_NAME}-latest-{{.OS}}-{{.Arch}}"

dist: build-cli-cross

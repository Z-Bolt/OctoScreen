# Packages content
PKG_OS = darwin linux
PKG_ARCH = amd64

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOINSTALL = $(GOCMD) install -v
GOGET = $(GOCMD) get -v -t
GOTEST = $(GOCMD) test -v

# Environment
BUILD_PATH := $(shell pwd)/build
DOCKER_IMAGE_BUILD = mcuadros/octoprint-tft-build

DEBIAN_PACKAGES = JESSIE
STRETCH_NAME := stretch
STRETCH_IMAGE := golang:1.9-stretch
STRETCH_GO_TAGS := gtk_3_22

JESSIE_NAME := jessie
JESSIE_IMAGE := golang:1.8-jessie
JESSIE_GO_TAGS := gtk_3_14


# Build information
GIT_COMMIT = $(shell git rev-parse HEAD | cut -c1-7)
GIT_DIRTY = $(shell test -n "`git status --porcelain`" && echo "-dirty" || true)
DEV_PREFIX := 0.dev
VERSION ?= $(DEV_PREFIX)-$(GIT_COMMIT)$(GIT_DIRTY)

ifneq ($(origin TRAVIS_TAG), undefined)
	VERSION := $(TRAVIS_TAG)
endif

# Package information
BINARY_NAME = octoprint-tft
INSTALL_FOLDER = /opt/octoprint-tft
PACKAGE_NAME =  ${BINARY_NAME}_${VERSION}_$(shell uname -m).deb

# we export the variable to allow envsubst, substitute the vars in the
# Dockerfiles
export

build-environment:
	mkdir -p ${BUILD_PATH}

build: | build-environment $(DEBIAN_PACKAGES)

$(DEBIAN_PACKAGES):
	docker build \
		--build-arg IMAGE=${${@}_IMAGE} \
		--build-arg GO_TAGS=${${@}_GO_TAGS} \
		-t ${DOCKER_IMAGE_BUILD}:${${@}_NAME} . \
		&& \
	docker run -it --rm \
		-v ${BUILD_PATH}/${${@}_NAME}:/build \
		${DOCKER_IMAGE_BUILD}:${${@}_NAME} \
		make package-internal

build-internal:
	go build --tags ${GO_TAGS} -v -o /build/bin/${BINARY_NAME} main.go

package-internal: build-internal
	mkdir -p /build/deb/${INSTALL_FOLDER}/bin/; \
	cp -rf styles /build/deb/${INSTALL_FOLDER}; \
	cp -rf /build/bin/${BINARY_NAME} /build/deb/${INSTALL_FOLDER}/bin/; \
	cp -rf etc/DEBIAN /build/deb/; \
	envsubst < etc/DEBIAN/control > /build/deb/DEBIAN/control; \
	dpkg-deb --build /build/deb/ /build/${PACKAGE_NAME}

clean:
	rm -rf ${BUILD_PATH}

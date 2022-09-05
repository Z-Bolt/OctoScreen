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
WORKDIR := $(shell pwd)
BUILD_PATH := $(WORKDIR)/build
GOCACHE_PATH = $(WORKDIR)/gocache
DOCKER_IMAGE_BUILD = mcuadros/octoprint-tft-build

DEBIAN_PACKAGES = BUSTER
ARCH = armhf
# ARCH = amd64

BULLSEYE_NAME := bullseye
BULLSEYE_IMAGE := golang:1.19-bullseye
BULLSEYE_GO_TAGS := "gtk_3_24 glib_deprecated glib_2_66"

# Buster's gtk 3.24.5 doesn't work with gtk_3_24 tag
# Using gtk_3_22 produces some deprecation warnings, but it compiles
# More info: https://github.com/gotk3/gotk3/issues/874
BUSTER_NAME := buster
BUSTER_IMAGE := golang:1.19-buster
BUSTER_GO_TAGS := "gtk_3_22 glib_deprecated glib_2_58 pango_1_42"

STRETCH_NAME := stretch
STRETCH_IMAGE := golang:1.19rc1-stretch
STRETCH_GO_TAGS := "gtk_3_22 glib_deprecated glib_2_50 cairo_1_14 pango_1_40"

# Build information
#GIT_COMMIT = $(shell git rev-parse HEAD | cut -c1-7)
VERSION := 2.8.0
BUILD_DATE ?= $(shell date --utc +%Y%m%d-%H:%M:%S)
#BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

#ifneq ($(BRANCH), master)
#	VERSION := $(shell echo $(BRANCH)| sed -e 's/v//g')
#endif

# Package information
PACKAGE_NAME = octoscreen

# we export the variable to allow envsubst, substitute the vars in the
# Dockerfiles
export

build-environment:
	mkdir -p ${BUILD_PATH}
	mkdir -p ${GOCACHE_PATH}

build: | build-environment $(DEBIAN_PACKAGES)

$(DEBIAN_PACKAGES):
	docker build \
		--build-arg IMAGE=${${@}_IMAGE} \
		--build-arg TARGET_ARCH=${ARCH} \
		--build-arg GO_TAGS=${${@}_GO_TAGS} \
		-t ${DOCKER_IMAGE_BUILD}:${${@}_NAME}-${ARCH} . \
		&& \
	docker run --rm \
		-e TARGET_ARCH=${ARCH} \
		-v ${BUILD_PATH}/${${@}_NAME}-${ARCH}:/build \
		-v ${GOCACHE_PATH}/${${@}_NAME}-${ARCH}:/gocache \
		${DOCKER_IMAGE_BUILD}:${${@}_NAME}-${ARCH} \
		make build-internal

build-internal: prepare-internal
	#go build --tags ${GO_TAGS} -v -o /build/bin/${BINARY_NAME} main.go
	cd $(WORKDIR); \
	GOCACHE=/gocache debuild --prepend-path=/usr/local/go/bin/ --preserve-env -us -uc -a${TARGET_ARCH}; \
	cp ../*.deb /build/;

prepare-internal:
	dch --create -v $(VERSION)-1 --package $(PACKAGE_NAME) --controlmaint empty; \
	cd $(WORKDIR)/..; \
	tar -czf octoscreen_$(VERSION).orig.tar.gz --exclude-vcs OctoScreen

clean:
	rm -rf ${BUILD_PATH}

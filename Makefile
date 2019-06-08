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
DOCKER_IMAGE_BUILD = mcuadros/octoprint-tft-build

DEBIAN_PACKAGES = STRETCH
STRETCH_NAME := stretch
STRETCH_IMAGE := golang:1.9-stretch
STRETCH_GO_TAGS := gtk_3_22

JESSIE_NAME := jessie
JESSIE_IMAGE := golang:1.8-jessie
JESSIE_GO_TAGS := gtk_3_14


# Build information
#GIT_COMMIT = $(shell git rev-parse HEAD | cut -c1-7)
#DEV_PREFIX := 1.0
VERSION := 1.3
BUILD_DATE ?= $(shell date --utc +%Y%m%d-%H:%M:%S)
#BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

#ifneq ($(BRANCH), master)
#	VERSION := $(shell echo $(BRANCH)| sed -e 's/v//g')
#endif

# Package information
PACKAGE_NAME = octoprint-tft

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
		make build-internal

build-internal: prepare-internal
	#go build --tags ${GO_TAGS} -v -o /build/bin/${BINARY_NAME} main.go
	cd $(WORKDIR); \
	debuild --prepend-path=/usr/local/go/bin/ --preserve-env -us -uc; \
	cp ../*.deb /build/;

prepare-internal:
	dch --create -v $(VERSION) --package $(PACKAGE_NAME) empty; \
	cd $(WORKDIR)/..; \
	tar -czf octoprint-tft_$(VERSION).orig.tar.gz --exclude-vcs OctoPrint-TFT

clean:
	rm -rf ${BUILD_PATH}

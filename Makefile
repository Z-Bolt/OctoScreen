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

BUSTER_NAME := buster
BUSTER_IMAGE := golang:1.15-buster
BUSTER_GO_TAGS := gtk_3_24

STRETCH_NAME := stretch
STRETCH_IMAGE := golang:1.9-stretch
STRETCH_GO_TAGS := gtk_3_22

JESSIE_NAME := jessie
JESSIE_IMAGE := golang:1.8-jessie
JESSIE_GO_TAGS := gtk_3_14


# Build information
VERSION_EPOCH := 1
GIT_COMMIT != git rev-parse HEAD | cut -c1-7
VERSION != git name-rev --tags --name-only $(GIT_COMMIT) | sed -e 's/(^[^0-9]+|[^0-9.])//g'
BRANCH != git rev-parse --abbrev-ref HEAD

# If this isn't the master branch, add additional development version information
ifneq ($(BRANCH), master)
    # Add the Safe Branch Name
    VERSION := $(VERSION)~$(shell echo $(BRANCH) | sed -e 's/[^A-Za-z0-9.]//g')
    # Add commits ahead of master
    VERSION := $(VERSION)+$(shell git rev-list --count HEAD ^master)
    # Add the current commit SHA (abreviated)
    VERSION := $(VERSION):$(GIT_COMMIT)
    # Add dirty indicator
    ifneq ($(shell git status --short | wc -l), 0)
        VERSION := $(VERSION)-dirty
    endif
endif

# Package information
PACKAGE_NAME = octoscreen

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
	dch --create -v '$(VERSION_EPOCH):$(VERSION)-1' --package $(PACKAGE_NAME) empty; \
	cd $(WORKDIR)/..; \
	tar -czf 'octoscreen_$(VERSION).orig.tar.gz' --exclude-vcs --force-local OctoScreen

clean:
	rm -rf ${BUILD_PATH}

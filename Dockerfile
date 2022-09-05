ARG IMAGE
FROM ${IMAGE}

SHELL ["/bin/bash", "-c"]

ARG TARGET_ARCH

# Note: dh-systemd doesn't exist on bullseye anymore (it's been merged with debhelper)
RUN apt-get update && \
  apt-get install -y --no-install-recommends \
    git build-essential \
    debhelper devscripts fakeroot git-buildpackage dh-make dh-golang \
  && \
  (apt-get install -y --no-install-recommends dh-systemd || true) \
  && \
  if [[ `dpkg-architecture -q DEB_BUILD_ARCH` != "${TARGET_ARCH}" ]]; then \
    dpkg --add-architecture ${TARGET_ARCH} && \
    apt-get update && \
    apt-get install -y --no-install-recommends \
      crossbuild-essential-${TARGET_ARCH}; \
  fi && \
  apt-get install -y --no-install-recommends \
    libcairo2-dev:${TARGET_ARCH} \
    libgtk-3-dev:${TARGET_ARCH}

ARG GO_TAGS
ENV GO_TAGS=${GO_TAGS}

# We cache go get gtk, to speed up builds.
#RUN go get -tags ${GO_TAGS} -v github.com/gotk3/gotk3/gtk/...

ADD . /OctoScreen/
#RUN go get -tags ${GO_TAGS} -v ./...

WORKDIR /OctoScreen/

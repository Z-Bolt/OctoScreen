FROM mcuadros/golang-arm:1.9-jessie
RUN apt update
RUN apt install -y --no-install-recommends git build-essential libcairo2-dev libgtk-3-dev

RUN rm -rf /go-linux-arm-bootstrap
ENV GOROOT=/usr/local/go/
RUN go get -tags gtk_3_14 -v github.com/gotk3/gotk3/gtk/...


ADD . /go/src/github.com/mcuadros/OctoPrint-TFT/

RUN go get -v github.com/sirupsen/logrus/...
WORKDIR /go/src/github.com/mcuadros/OctoPrint-TFT/
RUN go get -tags gtk_3_14 -v ./...
RUN go build -tags gtk_3_14  -v .

FROM golang:1.8

ADD . /go/

WORKDIR /go/src/github.com/guggero/docker-wallet-control

RUN go get -v ./...
RUN go install -v ./...

WORKDIR /go/

EXPOSE 80

CMD ["bin/docker-wallet-control"]

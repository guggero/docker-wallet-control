FROM golang:1.8

RUN mkdir -p /go/src/github.com/guggero/docker-wallet-control
WORKDIR /go/src/github.com/guggero/docker-wallet-control
COPY . .

RUN go-wrapper download
RUN go-wrapper install

EXPOSE 80

CMD ["go-wrapper", "run"]

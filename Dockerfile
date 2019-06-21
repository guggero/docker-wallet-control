FROM golang:1.12-stretch as builder

# Install dependencies and build the binaries.
RUN apt-get install -y \
    git \
    make \
    gcc \
&&  git clone https://github.com/guggero/docker-wallet-control /go/src/github.com/guggero/docker-wallet-control \
&&  cd /go/src/github.com/guggero/docker-wallet-control \
&&  make \
&&  make install

# Start a new, final image.
FROM debian:stretch as final

# Add bash and ca-certs, for quality of life and SSL-related reasons.
RUN apt-get update && apt-get install -y \
    bash \
    ca-certificates \
&&  mkdir /wallet-control

# Copy the binaries from the builder image.
COPY --from=builder /go/bin/docker-wallet-control /wallet-control/

# Copy the static resources
COPY static /wallet-control/static

WORKDIR /wallet-control

EXPOSE 80 443

# Specify the start command and entrypoint.
CMD ["/wallet-control/docker-wallet-control"]

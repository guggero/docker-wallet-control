# docker-wallet-control
Wallet admin control layer for cryptocoin wallets running in a docker environment.

**This is highly experimental and should only be used by people who know what they do!**

Running this app requires to open the RPC ports of cryptocoin wallets.
This can be very dangerous if done incorrectly and you could get all your coins stolen!

**You have been warned!**

## Docker command

Start the wallet control container with the following command:

```bash
docker run \
  -d \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -v /some/dir/config/config.json:/go/config.json \
  --restart always \
  --name wallet-control \
  guggero/docker-wallet-control
```

## Configuration file

You should create a file called `config.json` that you can then mount into the container as seen above.

Example:
```json
{
  "rpcUser": "jsonrpcuser",
  "rpcPassword": "some_safe_password_be_careful_or_your_money_will_get_stolen_you_have_been_warned!",
  "wallets": [
    {
      "label": "litecoin",
      "containerName": "litecoin",
      "rpcPort": 9332,
      "allowedUsers": ["test"]
    }
  ],
  "useClientCertAuth": false,
  "serveTLS": true,
  "serverAddress": "",
  "serverPort": 8443,
  "users": [
    {
      "username": "test",
      "password": "37268335dd6931045bdcdf92623ff819a64244b53d0e746d438797349d4da578",
      "salt": "test"
    }
  ]
}
```

## SSL/TLS configuration

To configure TLS (if either `useClientCertAuth` or `serveTLS` is true),
you need to mount a directory containing the certs to /go/tls:

```bash
docker run \
  -d \
  ...
  -v /some/dir/tls:/go/tls \
  ...
  guggero/docker-wallet-control
```

This directory should contain the following files:

* **server.key:** The RSA private key for the certificate without a password set
* **server.pem:** The server certificate
* **cacert.pem:** Optional, the CA used if client certificate authentication is enabled

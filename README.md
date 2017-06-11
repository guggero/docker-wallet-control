# docker-wallet-control
Wallet admin control layer for cryptocoin wallets running in a docker environment


## Configuration file

You should create a file called `config.json` that you can then mount into the container with:

```bash
docker run \
  -d \
  -v /some/dir/config.json:/go/config.json \
  ...
  guggero/docker-wallet-control
```

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

To configure TLS (if either `useClientCertAuth` or `serveTLS` is true), you need to mount
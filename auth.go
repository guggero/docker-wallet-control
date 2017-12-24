package main

import (
    "strings"
    "encoding/base64"
    "github.com/guggero/docker-wallet-control/util"
    "github.com/guggero/docker-wallet-control/rpc"
    "github.com/guggero/docker-wallet-control/docker"
    "net/http"
    "errors"
    "fmt"
)

func getAuthenticatedUser(r *http.Request) (string) {
    s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
    if len(s) != 2 {
        return ""
    }

    b, err := base64.StdEncoding.DecodeString(s[1])
    if err != nil {
        return ""
    }

    pair := strings.SplitN(string(b), ":", 2)
    if len(pair) != 2 {
        return ""
    }

    return authenticateUser(pair[0], pair[1])
}

func authenticateUser(username string, password string) (string) {
    for _, user := range appConfig.User {
        if user.Username == username && util.HashPassword(password, user.Salt) == user.Password {
            return user.Username
        }
    }
    return ""
}

func getRPCClient(walletName string, r *http.Request) (*rpc.Client, error) {
    user := getAuthenticatedUser(r)

    if user == "" {
        return nil, errors.New("User not authenticated")
    }

    for _, wallet := range appConfig.Wallets {
        if wallet.ContainerName == walletName {

            if util.ArrayContains(wallet.AllowedUsers, user) {
                var url = fmt.Sprintf("http://%s:%d/", wallet.ContainerName, wallet.RPCPort)
                client := rpc.CreateClient(url, appConfig.RPCUser, appConfig.RPCPassword)
                return client, nil
            } else {
                return nil, errors.New("User is not allowed to interact with wallet")
            }
        }
    }
    return nil, errors.New("Wallet not found")
}

func getDockerclient(walletName string, r *http.Request) (*docker.Client, error) {
    user := getAuthenticatedUser(r)

    if user == "" {
        return nil, errors.New("User not authenticated")
    }

    for _, wallet := range appConfig.Wallets {
        if wallet.ContainerName == walletName {

            if util.ArrayContains(wallet.AllowedUsers, user) {
                return docker.CreateClient(wallet.ContainerName)
            } else {
                return nil, errors.New("User is not allowed to interact with wallet")
            }
        }
    }
    return nil, errors.New("Wallet not found")
}

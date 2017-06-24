package main

import (
    "fmt"
    "time"
    "github.com/guggero/docker-wallet-control/util"
    "github.com/guggero/docker-wallet-control/rpc"
)

var appConfig *util.Configuration

func main() {

    var err error
    appConfig, err = util.ReadConfiguration("config.json")
    if err != nil {
        panic(err)
    }

    checkWallets()
    runServer()
}

func checkWallets() {
    for _, wallet := range appConfig.Wallets {

        var url = fmt.Sprintf("http://%s:%d/", wallet.ContainerName, wallet.RPCPort)
        client := rpc.CreateClient(url, appConfig.RPCUser, appConfig.RPCPassword)

        var start = time.Now()
        info := client.GetInfo()
        if (info.Blocks <= 0) {
            fmt.Printf("Error checking wallet %s, got %d blocks!", wallet.ContainerName, info.Blocks)
        }
        var duration = time.Since(start)
        fmt.Printf("Finished checking %s in %s\n",
            wallet.ContainerName,
            duration)
    }
}

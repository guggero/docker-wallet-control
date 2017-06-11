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

    runBenchmark()
    runServer()
}

func runBenchmark() {
    for _, wallet := range appConfig.Wallets {

        var url = fmt.Sprintf("http://%s:%d/", wallet.ContainerName, wallet.RPCPort)
        client := rpc.CreateClient(url, appConfig.RPCUser, appConfig.RPCPassword)

        var start = time.Now()
        var num = 1000

        for i := 0; i < num; i++ {
            info := client.GetInfo()
            if (info.Blocks <= 0) {
                panic("No blocks!")
            }
        }
        var duration = time.Since(start)
        fmt.Printf("Finished measurement for %s in %s with %f ops/s\n",
            wallet.ContainerName,
            duration,
            float64(num) / (float64(duration.Nanoseconds()) / 1e9))
    }
}

package util

import (
    "encoding/json"
    "os"
)

type WalletConfig struct {
    Label         string     `json:"label"`
    ContainerName string     `json:"containerName"`
    RPCPort       uint16     `json:"rpcPort"`
    AllowedUsers  []string   `json:"allowedUsers"`
}

type UserConfig struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Salt     string `json:"salt"`
}

type Configuration struct {
    RPCUser           string            `json:"rpcuser"`
    RPCPassword       string            `json:"rpcpassword"`
    Wallets           []WalletConfig    `json:"wallets"`
    UseClientCertAuth bool              `json:"useClientCertAuth"`
    ServeTLS          bool              `json:"serveTLS"`
    ServerAddress     string            `json:"serverAddress"`
    ServerPort        uint16            `json:"serverPort"`
    User              []UserConfig      `json:"users"`
}

func ReadConfiguration(path string) (*Configuration, error) {
    var file *os.File
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }

    decoder := json.NewDecoder(file)
    configuration := Configuration{}
    if err := decoder.Decode(&configuration); err != nil {
        return nil, err
    }
    return &configuration, nil
}
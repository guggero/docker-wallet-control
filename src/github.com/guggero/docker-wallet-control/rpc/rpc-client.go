package rpc

import (
    "github.com/ybbus/jsonrpc"
    "github.com/guggero/docker-wallet-control/util"
    "errors"
)

type Client struct {
    Url       string
    RPCClient *jsonrpc.RPCClient
}

type GetInfo struct {
    Version         float64 `json:"version"`
    ProtocolVersion float64 `json:"protocolversion"`
    WalletVersion   float64 `json:"walletversion"`
    Balance         float64 `json:"balance"`
    Blocks          float64 `json:"blocks"`
    Connections     float64 `json:"connections"`
    Difficulty      float64 `json:"difficulty"`
    Testnet         bool    `json:"testnet"`
    KeyPoolOldest   float64 `json:"keypoololdest"`
    KeyPoolSize     float64 `json:"keypoolsize"`
    PayTxFee        float64 `json:"paytxfee"`
    RelayFee        float64 `json:"relayfee"`
    Errors          string  `json:"errors"`
}

type GetWalletInfo struct {
    WalletVersion float64 `json:"walletversion"`
    Balance       float64 `json:"balance"`
    TxCount       float64 `json:"txcount"`
    KeyPoolOldest float64 `json:"keypoololdest"`
    KeyPoolSize   float64 `json:"keypoolsize"`
}

type Account struct {
    Name         string         `json:"name"`
    Addresses    []string       `json:"addresses"`
    Transactions []Transaction  `json:"transactions"`
    Balance      float64        `json:"balance"`
}

type Summary struct {
    Label              string       `json:"label"`
    WalletType         string       `json:"wallettype"`
    ContainerName      string       `json:"containername"`
    TxCount            float64      `json:"txcount"`
    Balance            float64      `json:"balance"`
    UnconfirmedBalance float64      `json:"unconfirmedbalance"`
    Blocks             float64      `json:"blocks"`
    Difficulty         float64      `json:"difficulty"`
    Testnet            bool         `json:"testnet"`
    KeyPoolSize        float64      `json:"keypoolsize"`
    Accounts           []Account    `json:"accounts"`
    Errors             string       `json:"errors"`
    MasternodeStatus   interface{}  `json:"masternodeStatus"`
    Logs               []string     `json:"logs"`
}

type Transaction struct {
    Account       string    `json:"account"`
    Address       string    `json:"address"`
    Category      string    `json:"category"`
    Amount        float64   `json:"amount"`
    Label         string    `json:"label"`
    Vout          float64   `json:"vout"`
    Confirmations float64   `json:"confirmations"`
    BlockHash     string    `json:"blockhash"`
    BlockIndex    float64   `json:"blockindex"`
    BlockTime     float64   `json:"blocktime"`
    TransactionId string    `json:"txid"`
    Time          float64   `json:"time"`
    TimeReceived  float64   `json:"timereceived"`
}

func CreateClient(url string, user string, password string) (*Client) {
    client := Client{
        Url: url,
        RPCClient: jsonrpc.NewRPCClient(url),
    }
    client.RPCClient.SetBasicAuth(user, password)
    return &client
}

func (client *Client) GetSummary(hostname string, walletType string, label string) (Summary) {
    info := client.GetInfo()
    walletinfo := client.GetWalletInfo()
    accountmap := client.ListAccounts()

    summary := Summary{
        Label: label,
        WalletType: walletType,
        ContainerName: hostname,
        TxCount: walletinfo.TxCount,
        Balance: info.Balance,
        UnconfirmedBalance: client.GetUnconfirmedBalance(),
        Blocks: info.Blocks,
        Difficulty: info.Difficulty,
        Testnet: info.Testnet,
        KeyPoolSize: info.KeyPoolSize,
        Errors: info.Errors,
        Accounts: make([]Account, len(accountmap)),
        MasternodeStatus: client.Masternode("status"),
    }

    var idx = 0
    for key, value := range accountmap {
        summary.Accounts[idx] = Account{
            Name: key,
            Balance: value,
            Addresses: client.GetAddressesByAccount(key),
            Transactions: client.ListTransactions(key),
        }
        idx++
    }

    return summary
}

func (client *Client) GetInfo() (*GetInfo) {
    var info GetInfo
    client.doCall("getinfo", &info)
    return &info
}

func (client *Client) GetWalletInfo() (*GetWalletInfo) {
    var info GetWalletInfo
    client.doCall("getwalletinfo", &info)
    return &info
}

func (client *Client) GetUnconfirmedBalance() (float64) {
    var result float64
    client.doCall("getunconfirmedbalance", &result)
    return result
}

func (client *Client) GetAddressesByAccount(account string) ([]string) {
    var result []string
    client.doCall("getaddressesbyaccount", &result, account)
    return result
}

func (client *Client) ListAccounts() (map[string]float64) {
    var result map[string]float64
    client.doCall("listaccounts", &result)
    return result
}

func (client *Client) ListTransactions(account string) ([]Transaction) {
    var result []Transaction
    client.doCall("listtransactions", &result, account)
    return result
}

func (client *Client) GetAccountAddress(account string) (string) {
    var result string
    client.doCall("getaccountaddress", &result, account)
    return result
}

func (client *Client) SendFrom(transaction *Transaction) (string) {
    var result string
    client.doCall("sendfrom", &result, transaction.Account, transaction.Address, transaction.Amount)
    return result
}

func (client *Client) Masternode(command string) (interface{}) {
    var result interface{}
    client.doCall("masternode", &result, command)
    return &result
}

func (client *Client) doCall(method string, result interface{}, params... interface{}) {
    response, err := client.RPCClient.Call(method, params...)
    if err != nil {
        util.LogError(errors.New("Could not connect to RPC server. Make sure the port " +
            "is correct, the server is running and username/password are correct."))
        util.LogError(err)
        return
    }
    if err := response.GetObject(result); err != nil {
        util.LogError(err)
        return
    }
    if result == nil {
        panic("Got nil result for url " + client.Url + " and method " + method)
    }
}


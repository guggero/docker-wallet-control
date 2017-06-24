package main

import (
    "net/http"
    "github.com/gorilla/mux"
    "github.com/guggero/docker-wallet-control/rpc"
    "github.com/guggero/docker-wallet-control/docker"
    "strconv"
)

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type SummaryResponse struct {
    Summaries []rpc.Summary             `json:"summaries"`
    UiData    map[string]*interface{}   `json:"uiData"`
}

type Routes []Route

var routes = Routes{
    Route{
        "Summary",
        "GET",
        "/summary",
        routeShowSummary,
    },
    Route{
        "AccountAddress",
        "GET",
        "/{wallet}/account/{account}",
        routeAccountAddress,
    },
    Route{
        "SendFrom",
        "POST",
        "/{wallet}/sendfrom",
        routeSendFrom,
    },
    Route{
        "SendFrom",
        "POST",
        "/{wallet}/sendfrom",
        routeSendFrom,
    },
    Route{
        "MasternodeStatus",
        "GET",
        "/{wallet}/masternode/status",
        routeMasternodeStatus,
    },
    Route{
        "MasternodeList",
        "GET",
        "/{wallet}/masternode/list",
        routeMasternodeList,
    },
    Route{
        "Logs",
        "GET",
        "/{wallet}/logs",
        routeLogs,
    },
    Route{
        "Health",
        "GET",
        "/{wallet}/health",
        routeHealth,
    },
    Route{
        "Restart",
        "GET",
        "/{wallet}/restart",
        routeRestart,
    },
}

func routeShowSummary(w http.ResponseWriter, r *http.Request) {
    var response = &SummaryResponse{
        UiData: appConfig.UiData,
    }
    for _, wallet := range appConfig.Wallets {
        client, err := getRPCClient(wallet.ContainerName, r)
        if err == nil && client != nil {
            response.Summaries = append(response.Summaries, client.GetSummary(wallet.ContainerName, wallet.Type, wallet.Label))
        }
    }
    writeJsonResponse(w, http.StatusOK, response)
}

func routeAccountAddress(w http.ResponseWriter, r *http.Request) {
    var vars = mux.Vars(r)
    client, err := getRPCClient(vars["wallet"], r)
    if err != nil {
        handleError(w, err)
        return
    }
    writeJsonResponse(w, http.StatusOK, client.GetAccountAddress(vars["account"]))
}

func routeSendFrom(w http.ResponseWriter, r *http.Request) {
    var (
        transaction rpc.Transaction
        client *rpc.Client
        err error
        vars = mux.Vars(r)
    )

    if err = readJsonBody(w, r, &transaction); err != nil {
        handleError(w, err)
        return
    }

    if client, err = getRPCClient(vars["wallet"], r); err != nil {
        handleError(w, err)
        return
    }

    writeJsonResponse(w, http.StatusOK, client.SendFrom(&transaction))
}

func routeMasternodeStatus(w http.ResponseWriter, r *http.Request) {
    var (
        client *rpc.Client
        err error
        vars = mux.Vars(r)
    )

    if client, err = getRPCClient(vars["wallet"], r); err != nil {
        handleError(w, err)
        return
    }

    writeJsonResponse(w, http.StatusOK, client.Masternode("status"))
}

func routeMasternodeList(w http.ResponseWriter, r *http.Request) {
    var (
        client *rpc.Client
        err error
        vars = mux.Vars(r)
    )

    if client, err = getRPCClient(vars["wallet"], r); err != nil {
        handleError(w, err)
        return
    }

    writeJsonResponse(w, http.StatusOK, client.Masternode("list"))
}

func routeLogs(w http.ResponseWriter, r *http.Request) {
    var (
        vars = mux.Vars(r)
        client *docker.Client
        numLines = 500
        lines []string
        err error
    )

    if client, err = getDockerclient(vars["wallet"], r); err != nil {
        handleError(w, err)
        return
    }

    if lines, err := strconv.Atoi(r.URL.Query().Get("lines")); err == nil {
        numLines = lines
    }

    if numLines > 10000 {
        numLines = 10000
    }

    if lines, err = client.GetLogs(numLines); err != nil {
        handleError(w, err)
        return
    }
    writeJsonResponse(w, http.StatusOK, lines)
}

func routeHealth(w http.ResponseWriter, r *http.Request) {
    var (
        vars = mux.Vars(r)
        client *docker.Client
        containerWrapper *docker.ContainerWrapper
        err error
    )

    if client, err = getDockerclient(vars["wallet"], r); err != nil {
        handleError(w, err)
        return
    }

    if containerWrapper, err = client.InspectContainer(); err != nil {
        handleError(w, err)
        return
    }

    writeJsonResponse(w, http.StatusOK, containerWrapper.Container.State)
}

func routeRestart(w http.ResponseWriter, r *http.Request) {
    var (
        vars = mux.Vars(r)
        client *docker.Client
        timeout uint = 10
        err error
    )

    if client, err = getDockerclient(vars["wallet"], r); err != nil {
        handleError(w, err)
        return
    }

    if timeoutParam, err := strconv.ParseUint(r.URL.Query().Get("timeout"), 10, 32); err == nil {
        timeout = uint(timeoutParam)
    }

    if err = client.Restart(timeout); err != nil {
        handleError(w, err)
        return
    }

    writeJsonResponse(w, http.StatusOK, "")
}
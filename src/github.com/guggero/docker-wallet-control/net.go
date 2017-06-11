package main

import (
    "net/http"
    "io/ioutil"
    "io"
    "encoding/json"
    "github.com/guggero/docker-wallet-control/util"
)

func readJsonBody(w http.ResponseWriter, r *http.Request, result interface{}) (error) {
    var (
        body []byte
        err error
    )
    if body, err = ioutil.ReadAll(io.LimitReader(r.Body, 1048576)); err != nil {
        return err
    }
    if err := r.Body.Close(); err != nil {
        return err
    }
    if err := json.Unmarshal(body, result); err != nil {
        return err
    }
    return nil
}

func handleError(w http.ResponseWriter, err error) {
    util.LogError(err)
    setJsonHeaders(w, http.StatusInternalServerError)
    if err := json.NewEncoder(w).Encode(err); err != nil {
        util.LogError(err)
        panic(err)
    }
}

func setJsonHeaders(w http.ResponseWriter, statusCode int) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(statusCode)
}

func writeJsonResponse(w http.ResponseWriter, statusCode int, entity interface{}) {
    setJsonHeaders(w, statusCode)
    if err := json.NewEncoder(w).Encode(entity); err != nil {
        handleError(w, err)
        return
    }
}

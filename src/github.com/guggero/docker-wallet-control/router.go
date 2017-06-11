package main

import (
    "github.com/gorilla/mux"
    "net/http"
)

func NewRouter() *mux.Router {

    router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
        var handler http.Handler

        handler = requestHandler(route.HandlerFunc, route.Name)

        router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
    }

    router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

    return router
}

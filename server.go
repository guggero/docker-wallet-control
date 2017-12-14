package main

import (
  "io/ioutil"
  "log"
  "crypto/x509"
  "crypto/tls"
  "net/http"
  "time"
  "fmt"
)

func runServer() {
  listenAddress := fmt.Sprintf("%s:%d", appConfig.ServerAddress, appConfig.ServerPort)

  if appConfig.ServeTLS {
    var tlsConfig *tls.Config

    if appConfig.UseClientCertAuth {
      certBytes, err := ioutil.ReadFile("tls/cacert.pem")
      if err != nil {
        log.Fatalln("Unable to read server.pem", err)
      }

      clientCertPool := x509.NewCertPool()
      if ok := clientCertPool.AppendCertsFromPEM(certBytes); !ok {
        log.Fatalln("Unable to add certificate to certificate pool")
      }

      tlsConfig = &tls.Config{
        ClientAuth:               tls.RequireAndVerifyClientCert,
        ClientCAs:                clientCertPool,
        PreferServerCipherSuites: true,
        MinVersion:               tls.VersionTLS12,
      }
      tlsConfig.BuildNameToCertificate()
    }

    server := &http.Server{
      Addr:      listenAddress,
      TLSConfig: tlsConfig,
      Handler:   NewRouter(),
    }
    log.Fatal(server.ListenAndServeTLS("tls/server.pem", "tls/server.key"))
  } else {
    log.Fatal(http.ListenAndServe(listenAddress, NewRouter()))
  }
}

func requestHandler(inner http.Handler, name string) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    start := time.Now()

    authUser := getAuthenticatedUser(r)
    authOk := authUser != ""

    if authOk && appConfig.UseClientCertAuth {
      if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 && len(r.TLS.PeerCertificates[0].EmailAddresses) > 0 {
        certUser := r.TLS.PeerCertificates[0].EmailAddresses[0]
        if certUser != authUser {
          log.Printf("Using client cert auth and basic auth, "+
            "but users don't match. Basic auth: %s, client cert auth: %s",
            authUser,
            certUser)
          authOk = false
        }
      } else {
        authOk = false
      }
    }

    if authOk {
      inner.ServeHTTP(w, r)
    } else {
      authUser = "-"
      w.Header().Set("WWW-Authenticate", "Basic realm=wallet-control")
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.WriteHeader(http.StatusUnauthorized)
      w.Write([]byte("Unauthorized\n"))
    }

    log.Printf(
      "[%v]\t\t%s\t%s\t%s\t%s",
      authUser,
      r.Method,
      r.RequestURI,
      name,
      time.Since(start),
    )
  })
}

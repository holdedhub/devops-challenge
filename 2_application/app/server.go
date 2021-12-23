package main

import (
  "fmt"
  "log"
  "net/http"
  "time"
)

func appHandler(w http.ResponseWriter, r *http.Request) {

  fmt.Println(time.Now(), "Hello from my new fresh server")

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/healthz" {
    http.Error(w, "404 not found.", http.StatusNotFound)
    return
  }

  if r.Method != "GET" {
    http.Error(w, "Method is not supported.", http.StatusNotFound)
    return
  }

  fmt.Fprintf(w, "ok")
}

func main() {
  http.HandleFunc("/", appHandler)
  http.HandleFunc("/healthz", healthHandler)

  log.Println("Started, serving on port 8080")
  err := http.ListenAndServe(":8080", nil)

  if err != nil {
    log.Fatal(err.Error())
  }
}

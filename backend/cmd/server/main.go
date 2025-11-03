package main

import (
    "log"
    "net/http"
    "os"
)

func main() {
    addr := ":8080"
    if p := os.Getenv("PORT"); p != "" {
        addr = ":" + p
    }

    http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    })

    log.Printf("Daily-Spark backend listening on %s", addr)
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatal(err)
    }
}
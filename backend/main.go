package main

import (
    "log"
    "net/http"

    "cbdc-backend/router"
)

func main() {
    r := router.SetupRouter()

    // Basic server configuration
    srv := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }

    log.Println("Starting CBDC backend on :8080")
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("server error: %v", err)
    }
}



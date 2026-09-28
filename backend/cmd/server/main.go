package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AndresVelezR/calculator---react-go/backend/internal/httpapi"
)

func main() {
	allowedOrigin := os.Getenv("CORS_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:5173"
	}
	server := &http.Server{
		Addr:              ":8080",
		Handler:           httpapi.CORS(httpapi.NewRouter(), allowedOrigin),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("server listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

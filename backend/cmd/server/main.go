package main

import (
	"log"
	"net/http"
	"os"

	"mirrorapp/backend/internal/server"
)

func main() {
	addr := getenv("ADDR", ":8080")

	app := server.New()
	handler := app.Routes()

	if dist := os.Getenv("FRONTEND_DIST"); dist != "" {
		handler = app.WithStatic(handler, dist)
	}

	log.Printf("Mirror backend listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

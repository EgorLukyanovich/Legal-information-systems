package main

import (
	"log"
	"net/http"
	"time"

	"github.com/egor_lukyanovich/legal-information-systems/backend/pkg/app"
	"github.com/egor_lukyanovich/legal-information-systems/backend/pkg/routing"
)

func main() {
	storage, port, err := app.InitDB()
	if err != nil {
		log.Fatalf("DB initialization failed: %v", err)
	}

	defer storage.DB.Close()

	router := routing.NewRouter(*storage)

	server := &http.Server{
		Handler:           router,
		Addr:              port,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on :%s", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to listen to server: %v", err)
	}
}

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/egor_lukyanovich/legal-information-systems/backend/internal/handlers"
	"github.com/egor_lukyanovich/legal-information-systems/backend/pkg/app"
	"github.com/egor_lukyanovich/legal-information-systems/backend/pkg/routing"
)

func main() {
	storage, port, err := app.InitDB()
	if err != nil {
		log.Fatalf("DB initialization failed: %v", err)
	}

	defer storage.DB.Close()

	userHandler := handlers.NewUserHandlers(storage.Queries)
	siteHandler := handlers.NewSiteHandlers(storage.Queries)
	router := routing.NewRouter(*storage)

	router.Post("/user/create", userHandler.CreateUser)
	router.Post("/user/auth", userHandler.UserAuthHandler)
	router.Post("/answers", handlers.AuthMiddleware(siteHandler.SubmitTestAnswers))
	router.Post("/theory/create", handlers.AuthMiddleware(siteHandler.CreateTheory))
	router.Post("/exemple/create", handlers.AuthMiddleware(siteHandler.CreateExample))

	router.Get("/user/profile", handlers.AuthMiddleware(userHandler.GetUserProfile))

	router.Get("/theory/get", handlers.AuthMiddleware(siteHandler.GetTheories))
	router.Get("/example/get", handlers.AuthMiddleware(siteHandler.GetExamples))
	router.Get("/questions/get", handlers.AuthMiddleware(siteHandler.GetQuestions))

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

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/egor_lukyanovich/legal-information-systems/backend/internal/handlers"
	"github.com/egor_lukyanovich/legal-information-systems/backend/pkg/app"
	"github.com/egor_lukyanovich/legal-information-systems/backend/pkg/routing"
	"github.com/go-chi/chi/v5"
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

	router.Group(func(r chi.Router) {
		r.Use(handlers.AuthMiddleware)

		r.Post("/answer", siteHandler.SubmitTestAnswers)
		r.Post("/theories", siteHandler.CreateTheory)
		r.Post("/examples", siteHandler.CreateExample)
		r.Post("/tests", siteHandler.CreateTest)

		r.Get("/user/profile", userHandler.GetUserProfile)
		r.Get("/theories", siteHandler.GetTheories)
		r.Get("/tests", siteHandler.GetTests)
		r.Get("/examples", siteHandler.GetExamples)
		r.Get("/questions", siteHandler.GetQuestions)
	})
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

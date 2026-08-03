package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"devops-monitoring-platform/backend/internal/db"
	"devops-monitoring-platform/backend/internal/handler"
	"devops-monitoring-platform/backend/internal/repository"
	"devops-monitoring-platform/backend/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	database, err := db.Connect()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer database.Close()

	repo := repository.NewPostgresRepository(database)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	router := chi.NewRouter()
	h.RegisterRoutes(router)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("backend listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

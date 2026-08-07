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
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env.local")
	if err != nil {
		log.Println("No .env file found")
	}

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

	router := newRouter(h)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("backend listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func newRouter(h *handler.Handler) chi.Router {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"http://127.0.0.1:5173", "http://localhost:5173", "http://127.0.0.1:80", "http://localhost:80"},
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	h.RegisterRoutes(router)

	return router
}

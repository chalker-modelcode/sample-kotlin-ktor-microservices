package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sample-kotlin-ktor-microservices/internal/config"
	"github.com/sample-kotlin-ktor-microservices/internal/handler"
	"github.com/sample-kotlin-ktor-microservices/internal/repository"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize repositories
	accountRepo := repository.NewInMemoryAccountRepository()

	// Initialize handlers
	accountHandler := handler.NewAccountHandler(accountRepo)

	// Set up router
	r := newRouter(accountHandler)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

// newRouter creates and configures the Chi router with all routes and middleware.
func newRouter(accountHandler *handler.AccountHandler) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// Health
	r.Get("/health", handler.HealthCheck)

	// Metrics
	r.Handle("/metrics", promhttp.Handler())

	// Account routes
	r.Route("/accounts", func(r chi.Router) {
		r.Post("/", accountHandler.CreateAccount)
		r.Get("/", accountHandler.ListAccounts)
		r.Get("/{id}", accountHandler.GetAccountByID)
		r.Get("/customer/{customerId}", accountHandler.GetAccountsByCustomerID)
	})

	return r
}

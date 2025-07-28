package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "effective_mobile/docs"
	"effective_mobile/src/_core/config"
	"effective_mobile/src/_core/db"
	"effective_mobile/src/subscriptions"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func enableCORS(router *mux.Router) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		router.ServeHTTP(w, r)
	})
}

// @title Subscription Service API
// @version 1.0
// @description API for managing user subscriptions

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:4000
// @BasePath /api
// @schemes http
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	gormDB, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.RunMigrations("up"); err != nil {
		log.Fatal("Migration failed: ", err)
	}

	subscriptionRepo := subscriptions.NewSubscriptionRepo(gormDB)
	subscriptionService := subscriptions.NewSubscriptionService(subscriptionRepo)
	subscriptionController := subscriptions.NewSubscriptionController(subscriptionService)

	// ROUTERS
	r := mux.NewRouter()
	corsRouter := enableCORS(r)
	api := r.PathPrefix("/api").Subrouter()
	subscriptionController.RegisterRoutes(api)

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DocExpansion("none"),
	))

	srv := &http.Server{
		Handler:      corsRouter,
		Addr:         fmt.Sprintf(":%d", cfg.API.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Starting server on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

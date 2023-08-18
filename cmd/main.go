package main

import (
	"fmt"
	"net/http"

	"github.com/gimlet-io/yaml-generator-app/cmd/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("App init..")

	err := godotenv.Load(".env")
	if err != nil {
		log.Warnf("could not load .env file, relying on env vars")
	}

	config, err := config.Environ()
	if err != nil {
		log.Fatalln("main: invalid configuration")
	}

	r := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"}, // Use this to allow specific origin hosts
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("config", config))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	r.Post("/", yamlGenerator)
	r.Post("/chart/{chart}", yamlGeneratorWithChart)

	err = http.ListenAndServe(":9000", r)
	log.Error(err)
}

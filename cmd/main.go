package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("App init..")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", yamlGenerator)
	r.Get("/ping", ping)

	http.ListenAndServe(":8080", r)
	err := http.ListenAndServe(":9000", r)
	log.Error(err)
}

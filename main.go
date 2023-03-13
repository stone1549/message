package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/stone1549/yapyapyap/message/service"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	flag.Parse()

	config, err := service.GetConfiguration()

	router := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Use(middleware.Timeout(time.Second * 30))
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Compress(5))
	router.Use(middleware.Recoverer)
	router.Use(middleware.AllowContentType("application/json"))

	configMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := context.WithValue(request.Context(), "config", config)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
	router.Use(configMiddleware)
	router.Use(service.JwtAuthMiddleware)

	repo, err := service.NewMessageRepository(config)

	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	repoMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := context.WithValue(request.Context(), "repo", repo)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
	router.Use(repoMiddleware)

	router.Route("/messages", func(r chi.Router) {
		r.With(service.GetMessagesMiddleware).Get("/", service.GetMessages)
		r.With(service.GetMessageMiddleware).Get("/{id}", service.GetMessage)
		r.With(service.AddMessageMiddleware).Put("/", service.AddMessage)
	})

	err = http.ListenAndServe(fmt.Sprintf(":%d", config.GetPort()), router)

	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}

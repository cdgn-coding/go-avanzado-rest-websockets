package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go-rest-websockets/handlers"
	"go-rest-websockets/middlewares"
	"go-rest-websockets/repository"
	"go-rest-websockets/server"
	"go-rest-websockets/websockets"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("cannot load configuration %v", err)
	}
	config := &server.Config{
		Port:        os.Getenv("PORT"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		DatabaseUrl: os.Getenv("DATABASE_URL"),
	}

	hub := websockets.NewHub()
	go hub.Run()
	s, err := server.NewServer(context.Background(), config, hub)
	if err != nil {
		log.Fatalf("cannot initialize server %v", err)
	}

	repo, err := repository.NewPostgresUserRepository(os.Getenv("DATABASE_URL"))
	authorization := server.NewAuthorization()
	bindRoutes := func(s server.Server, r *mux.Router) {
		r.Use(middlewares.CheckAuthMiddleware(s, authorization))
		r.Handle("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
		r.Handle("/signup", handlers.SignUpHandler(s, repo)).Methods(http.MethodPost)
		r.Handle("/login", handlers.LoginHandler(s, repo, authorization)).Methods(http.MethodPost)
		r.Handle("/me", handlers.MeHandler(s, repo, authorization)).Methods(http.MethodGet)
		r.Handle("/posts", handlers.InsertPostHandler(s, repo, authorization)).Methods(http.MethodPost)
		r.Handle("/posts/{id}", handlers.GetPostHandler(s, repo)).Methods(http.MethodGet)
		r.Handle("/posts/{id}", handlers.UpdatePostHandler(s, repo, authorization)).Methods(http.MethodPut)
		r.Handle("/posts/{id}", handlers.DeletePostHandler(s, repo, authorization)).Methods(http.MethodDelete)
		r.Handle("/posts", handlers.GetPaginatedPostsHandler(s, repo)).Methods(http.MethodGet)
		r.HandleFunc("/ws", s.Hub().HandleWebSocket)
	}

	s.Start(bindRoutes)
}

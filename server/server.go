package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go-rest-websockets/websockets"
	"log"
	"net/http"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseUrl string
}

type Server interface {
	Config() *Config
	Hub() *websockets.Hub
}

type Broker struct {
	config *Config
	router *mux.Router
	hub    *websockets.Hub
}

func NewServer(ctx context.Context, config *Config, hub *websockets.Hub) (*Broker, error) {
	if config.Port == "" {
		return nil, fmt.Errorf("port is required")
	}

	if config.JWTSecret == "" {
		return nil, fmt.Errorf("jwt secret is required")
	}

	if config.DatabaseUrl == "" {
		return nil, fmt.Errorf("database url is required")
	}

	return &Broker{config: config, router: mux.NewRouter(), hub: hub}, nil
}

func (b *Broker) Config() *Config {
	return b.config
}

type BinderFunc func(s Server, r *mux.Router)

func (b *Broker) Start(binder BinderFunc) {
	b.router = mux.NewRouter()
	binder(b, b.router)

	log.Print("Starting server on port", b.Config().Port)
	err := http.ListenAndServe(b.Config().Port, b.router)

	if err != nil {
		log.Fatalf("Error starting server %v", err)
	}
}

func (b *Broker) Hub() *websockets.Hub {
	return b.hub
}

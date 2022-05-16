package handlers

import (
	"encoding/json"
	"go-rest-websockets/server"
	"log"
	"net/http"
)

type HomeResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func HomeHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := HomeResponse{
			Message: "This is the Home Controller",
			Status:  true,
		}
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		err := encoder.Encode(response)
		if err != nil {
			log.Printf("error encoding response %v", err)
		}
	}
}

package handlers

import (
	"encoding/json"
	"github.com/segmentio/ksuid"
	"go-rest-websockets/models"
	"go-rest-websockets/repository"
	"go-rest-websockets/server"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

const HashCost = 8

type SignUpLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type MeResponse struct {
	Email string `json:"email"`
	Id    string `json:"id"`
}

type LoginResponse struct {
	UserId string `json:"userId"`
	Token  string `json:"token"`
}

func SignUpHandler(s server.Server, r repository.Repository) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		signUpRequest := SignUpLoginRequest{}
		err := json.NewDecoder(request.Body).Decode(&signUpRequest)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(signUpRequest.Password), HashCost)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		signUpUser := models.User{
			Id:       id.String(),
			Email:    signUpRequest.Email,
			Password: string(hashed),
		}
		err = r.InsertUser(request.Context(), &signUpUser)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
		err = json.NewEncoder(writer).Encode(SignUpResponse{
			Id:    id.String(),
			Email: signUpUser.Email,
		})
		if err != nil {
			log.Printf("error encoding response %v", err)
		}
	}
}

func LoginHandler(s server.Server, r repository.Repository, auth server.Authorization) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		loginRequest := SignUpLoginRequest{}
		err := json.NewDecoder(request.Body).Decode(&loginRequest)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := r.GetUserByEmail(request.Context(), loginRequest.Email)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		if user == nil {
			http.Error(writer, "invalid credentials", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
		if user == nil {
			http.Error(writer, "invalid credentials", http.StatusUnauthorized)
			return
		}

		tokenString, err := auth.SignToken(s.Config().JWTSecret, user.Id)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		loginResponse := LoginResponse{
			UserId: user.Id,
			Token:  tokenString,
		}

		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(writer).Encode(loginResponse)
		if err != nil {
			log.Printf("error while encoding %v", err)
		}
	}
}

func MeHandler(s server.Server, repository repository.Repository, auth server.Authorization) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		claims, err := auth.ParseAndVerifyToken(s.Config().JWTSecret, tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		user, err := repository.GetUserById(r.Context(), claims.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := MeResponse{
			Email: user.Email,
			Id:    user.Id,
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("error encoding %v", err)
		}
	}
}

package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
	"go-rest-websockets/models"
	"go-rest-websockets/repository"
	"go-rest-websockets/server"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type UpsertPostRequest struct {
	PostContent string `json:"postContent"`
}

func InsertPostHandler(s server.Server, repo repository.Repository, auth server.Authorization) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		insertPostRequest := UpsertPostRequest{}
		err := json.NewDecoder(r.Body).Decode(&insertPostRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		claims, err := auth.ParseAndVerifyToken(s.Config().JWTSecret, tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		post := models.Post{
			Id:          id.String(),
			PostContent: insertPostRequest.PostContent,
			UserId:      claims.UserId,
		}

		err = repo.InsertPost(r.Context(), &post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		message := models.WebsocketMessage{
			Type:    "Post_Created",
			Payload: post,
		}
		s.Hub().Broadcast(message, nil)

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			log.Printf("error encoding response %v", err)
		}
	}
}

func GetPostHandler(s server.Server, repository repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		post, err := repository.GetPostById(r.Context(), params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if post.Id == "" {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			log.Println("error encoding response")
		}
	}
}

func UpdatePostHandler(s server.Server, repository repository.Repository, auth server.Authorization) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		postId := params["id"]
		updatePostRequest := UpsertPostRequest{}
		err := json.NewDecoder(r.Body).Decode(&updatePostRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		claims, err := auth.ParseAndVerifyToken(s.Config().JWTSecret, tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		post := models.Post{
			Id:          postId,
			PostContent: updatePostRequest.PostContent,
			UserId:      claims.UserId,
		}
		err = repository.UpdatePost(r.Context(), &post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			log.Println("error encoding response")
		}
	}
}

func DeletePostHandler(s server.Server, repository repository.Repository, auth server.Authorization) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		postId := params["id"]
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		claims, err := auth.ParseAndVerifyToken(s.Config().JWTSecret, tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		post := models.Post{
			Id:     postId,
			UserId: claims.UserId,
		}
		err = repository.DeletePost(r.Context(), &post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func getIntQueryParam(r *http.Request, key string, defaultVal int) (int, error) {
	query := r.URL.Query()
	queryVals := query[key]
	var val int
	var err error

	if len(queryVals) != 1 {
		val = defaultVal
	} else {
		val, err = strconv.Atoi(queryVals[0])
	}
	return val, err
}

func GetPaginatedPostsHandler(s server.Server, repository repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, err := getIntQueryParam(r, "page", 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		posts, err := repository.GetPaginatedPosts(r.Context(), 10, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(posts)
		if err != nil {
			log.Println("error encoding response")
		}
	}
}

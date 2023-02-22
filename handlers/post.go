package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
	"platzi.com/go/rest-ws/models"
	"platzi.com/go/rest-ws/repository"
	"platzi.com/go/rest-ws/server"
)

type InsertPostRequest struct {
	Content string `json:"content"`
}

type PostResponse struct {
	Id      string `json:"id"`
	Content string `json:"content"`
}

func InsertPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			var postRequest = InsertPostRequest{}
			if err := json.NewDecoder(r.Body).Decode(&postRequest); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			id, err := ksuid.NewRandom()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			post := models.Post{
				Id:      id.String(),
				Content: postRequest.Content,
				UserId:  claims.UserId,
			}
			err = repository.InsertPost(r.Context(), &post)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(PostResponse{
				Id:      post.Id,
				Content: post.Content,
			})
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

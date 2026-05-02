package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var jwtKey []byte

func main() {

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	jwtKey = []byte(secret)

	InitDB()
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.With(JWTMiddleware).Get("/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		username := r.Context().Value(usernameKey).(string)

		json.NewEncoder(w).Encode(map[string]string{
			"username": username,
		})
	})

	r.Post("/register", RegisterHandler)
	r.Post("/login", LoginHandler)

	r.Route("/rooms", func(r chi.Router) {
		r.Use(JWTMiddleware)
		r.Post("/", CreateChatroomHandler)
		r.Get("/", GetChatroomsHandler)
		r.Patch("/{id}", UpdateChatroomHandler)
		r.Delete("/{id}", DeleteChatroomHandler)
		r.Post("/{id}/messages", CreateMessageHandler)
		r.Get("/{id}/messages", GetMessagesHandler)
		r.Delete("/{id}/messages/{msgId}", DeleteMessageHandler)
		r.Patch("/{id}/messages/{msgId}", UpdateMessageHandler)
	})

	r.Handle("/*", http.FileServer(http.Dir("./static")))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

package server

import (
	"chatapi/internal/auth"
	"chatapi/internal/chatroom"
	"chatapi/internal/message"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(
	authHandler *auth.Handler,
	chatroomHandler *chatroom.Handler,
	messageHandler *message.Handler,
) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)

	r.With(authHandler.Middleware).Get("/me", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		username := r.Context().Value(auth.UsernameKey).(string)
		json.NewEncoder(w).Encode(map[string]string{"username": username})
	})

	r.Route("/rooms", func(r chi.Router) {
		r.Use(authHandler.Middleware)
		r.Post("/", chatroomHandler.Create)
		r.Get("/", chatroomHandler.List)
		r.Delete("/{id}", chatroomHandler.Delete)
		r.Patch("/{id}", chatroomHandler.Update)
		r.Post("/{id}/messages", messageHandler.Create)
		r.Get("/{id}/messages", messageHandler.List)
		r.Delete("/{id}/messages/{msgId}", messageHandler.Delete)
		r.Patch("/{id}/messages/{msgId}", messageHandler.Update)

	})

	r.Handle("/*", http.FileServer(http.Dir("./static")))

	return r
}

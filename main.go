package main

import (
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
)

var users = make(map[string]User)
var chatrooms = make(map[int]Chatroom)
var nextRoomID = 1
var nextMessageID = 1
var jwtKey = []byte(os.Getenv("JWT_SECRET"))
var mu sync.Mutex

func main() {
	r := chi.NewRouter()

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

	http.ListenAndServe(":8080", r)
}

package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func main() {
	InitDB()
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

	r.Handle("/*", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":8080", r)
}

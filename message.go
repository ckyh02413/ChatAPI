package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Message struct {
	ID      int    `json:"id"`
	Creator string `json:"creator"`
	Content string `json:"content"`
}

func CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := r.Context().Value("username").(string)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid room id", http.StatusBadRequest)
		return
	}

	var input Message
	json.NewDecoder(r.Body).Decode(&input)

	if input.Content == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	chatroom, ok := getChatroom(w, id)

	if !ok {
		return
	}

	input.ID = nextMessageID
	nextMessageID++
	input.Creator = username

	chatroom.Messages = append(chatroom.Messages, input)
	chatrooms[id] = chatroom

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)

}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid room id", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	chatroom, ok := getChatroom(w, id)

	if !ok {
		return
	}

	json.NewEncoder(w).Encode(chatroom.Messages)
}

func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid room id", http.StatusBadRequest)
		return
	}

	messageID, err := strconv.Atoi(chi.URLParam(r, "msgId"))

	if err != nil {
		http.Error(w, "Invalid message id", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	chatroom, ok := getChatroom(w, id)

	if !ok {
		return
	}

	found := false

	for i, msg := range chatroom.Messages {
		if msg.ID == messageID {
			if msg.Creator != username && chatroom.Creator != username {
				http.Error(w, "Permission denied", http.StatusForbidden)
				return
			}
			chatroom.Messages = append(chatroom.Messages[:i], chatroom.Messages[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	chatrooms[id] = chatroom

}

func UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := r.Context().Value("username").(string)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid room id", http.StatusBadRequest)
		return
	}

	messageID, err := strconv.Atoi(chi.URLParam(r, "msgId"))
	if err != nil {
		http.Error(w, "Invalid message id", http.StatusBadRequest)
		return
	}

	var input Message
	json.NewDecoder(r.Body).Decode(&input)

	if input.Content == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	chatroom, ok := getChatroom(w, id)

	if !ok {
		return
	}

	found := false

	var updatedID int
	for i, msg := range chatroom.Messages {
		if msg.ID == messageID {
			if msg.Creator != username {
				http.Error(w, "Permission denied", http.StatusForbidden)
				return
			}
			chatroom.Messages[i].Content = input.Content
			updatedID = i
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	chatrooms[id] = chatroom
	json.NewEncoder(w).Encode(chatrooms[id].Messages[updatedID])
}

package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Chatroom struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Creator  string    `json:"creator"`
	Messages []Message `json:"messages"`
}

type ChatroomSummary struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getChatroom(w http.ResponseWriter, id int) (Chatroom, bool) {
	chatroom, exists := chatrooms[id]

	if !exists {
		http.Error(w, "Room not found", http.StatusNotFound)
		return Chatroom{}, false
	}

	return chatroom, true
}

func CreateChatroomHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := r.Context().Value("username").(string)

	var input Chatroom
	json.NewDecoder(r.Body).Decode(&input)

	if input.Name == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for _, cr := range chatrooms {
		if cr.Name == input.Name {
			http.Error(w, "Chatroom already exists", http.StatusConflict)
			return
		}
	}
	input.Creator = username
	input.ID = nextRoomID
	nextRoomID++
	chatrooms[input.ID] = input

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(ChatroomSummary{
		ID:   input.ID,
		Name: input.Name,
	})
}

func GetChatroomsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var summaries []ChatroomSummary

	mu.Lock()
	defer mu.Unlock()

	for _, cr := range chatrooms {
		summaries = append(summaries, ChatroomSummary{
			ID:   cr.ID,
			Name: cr.Name,
		})
	}
	json.NewEncoder(w).Encode(summaries)
}

func DeleteChatroomHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	chatroom, ok := getChatroom(w, id)

	if !ok {
		return
	}

	if chatroom.Creator != username {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}
	delete(chatrooms, id)
}

func UpdateChatroomHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := r.Context().Value("username").(string)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	var input Chatroom
	json.NewDecoder(r.Body).Decode(&input)

	if input.Name == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	chatroom, exists := chatrooms[id]

	if !exists {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	if chatroom.Creator != username {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	chatroom.Name = input.Name
	chatrooms[id] = chatroom

	json.NewEncoder(w).Encode(chatrooms[id])
}

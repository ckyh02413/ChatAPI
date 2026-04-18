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
	var chatroom Chatroom

	err := db.QueryRow(
		"SELECT id, name, creator FROM chatrooms WHERE id = $1",
		id,
	).Scan(&chatroom.ID, &chatroom.Name, &chatroom.Creator)

	if err != nil {
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

	var exists bool

	err := db.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM chatrooms WHERE name = $1)",
		input.Name,
	).Scan(&exists)

	if err != nil || exists {
		http.Error(w, "Chatroom already exists", http.StatusConflict)
		return
	}

	var id int
	err = db.QueryRow(
		"INSERT INTO chatrooms (name, creator) VALUES ($1, $2) RETURNING id",
		input.Name, username,
	).Scan(&id)

	if err != nil {
		http.Error(w, "Error creating chatroom", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(ChatroomSummary{
		ID:   id,
		Name: input.Name,
	})
}

func GetChatroomsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, name FROM chatrooms")

	if err != nil {
		http.Error(w, "Error fetching chatrooms", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var summaries []ChatroomSummary

	for rows.Next() {
		var summary ChatroomSummary
		err := rows.Scan(&summary.ID, &summary.Name)

		if err != nil {
			http.Error(w, "Error scanning chatroom", http.StatusInternalServerError)
			return
		}
		summaries = append(summaries, summary)
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

	chatroom, ok := getChatroom(w, id)

	if !ok {
		return
	}

	if chatroom.Creator != username {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	_, err = db.Exec(
		"DELETE FROM chatrooms WHERE id = $1",
		chatroom.ID,
	)

	if err != nil {
		http.Error(w, "Error deleting chatroom", http.StatusInternalServerError)
		return
	}
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

	chatroom, ok := getChatroom(w, id)

	if !ok {
		return
	}

	if chatroom.Creator != username {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	_, err = db.Exec(
		"UPDATE chatrooms SET name = $1 WHERE id = $2",
		input.Name, chatroom.ID,
	)

	if err != nil {
		http.Error(w, "Error updating chatroom", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ChatroomSummary{
		ID:   id,
		Name: input.Name,
	})
}

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
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)
	w.Header().Set("Content-Type", "application/json")

	username := r.Context().Value(usernameKey).(string)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid room id", http.StatusBadRequest)
		return
	}

	var input Message
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.Content == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	_, ok := getChatroom(w, id)

	if !ok {
		return
	}

	var messageID int

	err = db.QueryRow(
		"INSERT INTO messages (chatroom_id, creator, content) VALUES ($1, $2, $3) RETURNING id",
		id, username, input.Content,
	).Scan(&messageID)

	if err != nil {
		http.Error(w, "Error sending message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Message{
		ID:      messageID,
		Creator: username,
		Content: input.Content,
	})

}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid room id", http.StatusBadRequest)
		return
	}

	_, ok := getChatroom(w, id)

	if !ok {
		return
	}

	rows, err := db.Query(
		"SELECT id, creator, content FROM messages WHERE chatroom_id = $1",
		id,
	)

	if err != nil {
		http.Error(w, "Error fetching messages", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	messages := make([]Message, 0)

	for rows.Next() {
		var message Message
		err = rows.Scan(&message.ID, &message.Creator, &message.Content)

		if err != nil {
			http.Error(w, "Error scanning message", http.StatusInternalServerError)
			return
		}
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Error reading rows", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messages)
}

func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(usernameKey).(string)

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

	chatroom, ok := getChatroom(w, id)

	if !ok {
		return
	}

	var message Message

	err = db.QueryRow(
		"SELECT id, creator, content FROM messages WHERE id = $1",
		messageID,
	).Scan(&message.ID, &message.Creator, &message.Content)

	if err != nil {
		http.Error(w, "Error fetching message", http.StatusInternalServerError)
		return
	}

	if message.Creator != username && chatroom.Creator != username {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	_, err = db.Exec(
		"DELETE FROM messages WHERE id = $1",
		messageID,
	)

	if err != nil {
		http.Error(w, "Error deleting message", http.StatusInternalServerError)
		return
	}

}

func UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)
	w.Header().Set("Content-Type", "application/json")

	username := r.Context().Value(usernameKey).(string)

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
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.Content == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	_, ok := getChatroom(w, id)

	if !ok {
		return
	}

	var message Message

	err = db.QueryRow(
		"SELECT id, creator, content FROM messages WHERE id = $1",
		messageID,
	).Scan(&message.ID, &message.Creator, &message.Content)

	if err != nil {
		http.Error(w, "Error fetching message", http.StatusInternalServerError)
		return
	}

	if message.Creator != username {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	_, err = db.Exec(
		"UPDATE messages SET content = $1 WHERE id = $2",
		input.Content, messageID,
	)

	if err != nil {
		http.Error(w, "Error updating message", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(
		Message{
			ID:      message.ID,
			Creator: message.Creator,
			Content: input.Content,
		})
}

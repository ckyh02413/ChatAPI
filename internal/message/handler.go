package message

import (
	"chatapi/internal/auth"
	apperrors "chatapi/internal/errors"
	"chatapi/internal/validation"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)
	w.Header().Set("Content-Type", "application/json")

	username := r.Context().Value(auth.UsernameKey).(string)

	chatroomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid chatroom id", http.StatusBadRequest)
		return
	}

	var input CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validation.Validate(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := h.service.Create(chatroomID, username, input.Content)
	if err != nil {
		http.Error(w, err.Error(), apperrors.ErrorToStatus(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	chatroomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid room id", http.StatusBadRequest)
		return
	}

	messages, err := h.service.ListByChatroom(chatroomID)
	if err != nil {
		http.Error(w, err.Error(), apperrors.ErrorToStatus(err))
		return
	}

	json.NewEncoder(w).Encode(messages)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(auth.UsernameKey).(string)

	chatroomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid chatroom id", http.StatusBadRequest)
		return
	}

	messageID, err := strconv.Atoi(chi.URLParam(r, "msgId"))
	if err != nil {
		http.Error(w, "Invalid message id", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(chatroomID, messageID, username)
	if err != nil {
		http.Error(w, err.Error(), apperrors.ErrorToStatus(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)
	w.Header().Set("Content-Type", "application/json")

	username := r.Context().Value(auth.UsernameKey).(string)

	messageID, err := strconv.Atoi(chi.URLParam(r, "msgId"))
	if err != nil {
		http.Error(w, "Invalid message id", http.StatusBadRequest)
		return
	}

	var input UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validation.Validate(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := h.service.Update(messageID, username, input.Content)
	if err != nil {
		http.Error(w, err.Error(), apperrors.ErrorToStatus(err))
		return
	}

	json.NewEncoder(w).Encode(message)

}

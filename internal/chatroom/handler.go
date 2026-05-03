package chatroom

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

	var input CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validation.Validate(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chatroomSummary, err := h.service.Create(input.Name, username)
	if err != nil {
		http.Error(w, err.Error(), apperrors.ErrorToStatus(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chatroomSummary)

}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	chatroomSummaries, err := h.service.List()
	if err != nil {
		http.Error(w, err.Error(), apperrors.ErrorToStatus(err))
		return
	}

	json.NewEncoder(w).Encode(chatroomSummaries)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(auth.UsernameKey).(string)

	chatroomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(chatroomID, username)
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

	chatroomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
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

	chatroomSummary, err := h.service.Update(chatroomID, username, input.Name)
	if err != nil {
		http.Error(w, err.Error(), apperrors.ErrorToStatus(err))
		return
	}

	json.NewEncoder(w).Encode(chatroomSummary)
}

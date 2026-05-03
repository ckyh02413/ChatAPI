package auth

import (
	apperrors "chatapi/internal/errors"
	"chatapi/internal/validation"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const UsernameKey contextKey = "username"

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

type Handler struct {
	service *Service
	repo    *Repository
	jwtKey  []byte
}

func NewHandler(service *Service, repo *Repository, jwtKey []byte) *Handler {
	return &Handler{service: service, repo: repo, jwtKey: jwtKey}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)
	w.Header().Set("Content-Type", "application/json")

	var input RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validation.Validate(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.service.Register(input.Username, input.Mail, input.Password)
	if err != nil {
		http.Error(w, err.Error(), apperrors.ErrorToStatus(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)
	w.Header().Set("Content-Type", "application/json")

	var input LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validation.Validate(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(input.Username, input.Password)
	if err != nil {
		http.Error(w, err.Error(), apperrors.ErrorToStatus(err))
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (h *Handler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")

		if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
			http.Error(w, "Missing or malformed token", http.StatusUnauthorized)
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return h.jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		exists, err := h.repo.Exists(claims.Username)

		if err != nil || !exists {
			http.Error(w, "User no longer exists", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UsernameKey, claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

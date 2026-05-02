package auth

import (
	apperrors "chatapi/internal/errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	repo   *Repository
	jwtKey []byte
}

func NewService(repo *Repository, jwtKey []byte) *Service {
	return &Service{repo: repo, jwtKey: jwtKey}
}

func (s *Service) Register(username, mail, password string) (string, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return "", err
	}

	if err := s.repo.Create(username, mail, hashedPassword); err != nil {
		return "", apperrors.ErrAlreadyExists
	}

	return s.generateJWT(username)
}

func (s *Service) Login(username, password string) (string, error) {
	hashedPassword, err := s.repo.FindPasswordByUsername(username)
	if err != nil {
		return "", apperrors.ErrInvalidCredentials
	}

	if err := checkPassword(hashedPassword, password); err != nil {
		return "", apperrors.ErrInvalidCredentials
	}

	return s.generateJWT(username)
}

func (s *Service) generateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtKey)
}

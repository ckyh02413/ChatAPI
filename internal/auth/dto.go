package auth

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Mail     string `json:"mail" validate:"required,email,max=254"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

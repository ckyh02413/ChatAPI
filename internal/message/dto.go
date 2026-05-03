package message

type CreateRequest struct {
	Content string `json:"content" validate:"required,min=1,max=2000"`
}

type UpdateRequest struct {
	Content string `json:"content" validate:"required,min=1,max=2000"`
}

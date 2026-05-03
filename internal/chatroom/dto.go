package chatroom

type CreateRequest struct {
	Name string `json:"name" validate:"required,min=3,max=64"`
}

type UpdateRequest struct {
	Name string `json:"name" validate:"required,min=3,max=64"`
}

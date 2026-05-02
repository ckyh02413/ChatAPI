package chatroom

import "time"

type Chatroom struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Creator   string    `json:"creator"`
	CreatedAt time.Time `json:"created_at"`
}

type Summary struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

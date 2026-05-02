package message

import "time"

type Message struct {
	ID        int       `json:"id"`
	Creator   string    `json:"creator"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

package message

type Message struct {
	ID      int    `json:"id"`
	Creator string `json:"creator"`
	Content string `json:"content"`
}

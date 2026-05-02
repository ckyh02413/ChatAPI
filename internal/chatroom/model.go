package chatroom

type Chatroom struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Creator string `json:"creator"`
}

type Summary struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

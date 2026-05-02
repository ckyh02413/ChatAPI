package message

import (
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(chatroomID int, creator, content string) (int, time.Time, error) {
	var messageID int
	var createdAt time.Time
	err := r.db.QueryRow(
		"INSERT INTO messages (chatroom_id, creator, content) VALUES ($1, $2, $3) RETURNING id, created_at",
		chatroomID, creator, content,
	).Scan(&messageID, &createdAt)

	return messageID, createdAt, err
}

func (r *Repository) FindByID(messageID int) (Message, error) {
	var message Message
	err := r.db.QueryRow(
		"SELECT id, creator, content, created_at FROM messages WHERE id = $1",
		messageID,
	).Scan(&message.ID, &message.Creator, &message.Content, &message.CreatedAt)

	return message, err
}

func (r *Repository) ListByChatroom(chatroomID int) ([]Message, error) {
	rows, err := r.db.Query(
		"SELECT id, creator, content, created_at FROM messages WHERE chatroom_id = $1 ORDER BY created_at ASC",
		chatroomID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	messages := make([]Message, 0)
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.Creator, &message.Content, &message.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *Repository) UpdateContent(messageID int, newContent string) error {
	_, err := r.db.Exec(
		"UPDATE messages SET content = $1 WHERE id = $2",
		newContent, messageID,
	)

	return err
}

func (r *Repository) Delete(messageID int) error {
	_, err := r.db.Exec(
		"DELETE FROM messages WHERE id = $1",
		messageID,
	)

	return err
}

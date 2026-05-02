package chatroom

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(name, creator string) (int, error) {
	var chatroomID int
	err := r.db.QueryRow(
		"INSERT INTO chatrooms (name, creator) VALUES ($1, $2) RETURNING id",
		name, creator,
	).Scan(&chatroomID)

	return chatroomID, err
}

func (r *Repository) Exists(name string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM chatrooms WHERE name = $1)",
		name,
	).Scan(&exists)

	return exists, err
}

func (r *Repository) FindByID(chatroomID int) (Chatroom, error) {
	var chatroom Chatroom
	err := r.db.QueryRow(
		"SELECT id, name, creator FROM chatrooms WHERE id = $1",
		chatroomID,
	).Scan(&chatroom.ID, &chatroom.Name, &chatroom.Creator)

	return chatroom, err
}

func (r *Repository) List() ([]Summary, error) {
	rows, err := r.db.Query("SELECT id, name FROM chatrooms")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var summaries = make([]Summary, 0)
	for rows.Next() {
		var summary Summary
		if err := rows.Scan(&summary.ID, &summary.Name); err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return summaries, nil
}

func (r *Repository) UpdateName(chatroomID int, newName string) error {
	_, err := r.db.Exec(
		"UPDATE chatrooms SET name = $1 WHERE id = $2",
		newName, chatroomID,
	)

	return err
}

func (r *Repository) Delete(chatroomID int) error {
	_, err := r.db.Exec(
		"DELETE FROM chatrooms WHERE id = $1",
		chatroomID,
	)

	return err
}

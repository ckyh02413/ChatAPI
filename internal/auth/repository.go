package auth

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(username, mail, hashedPassword string) error {
	_, err := r.db.Exec(
		"INSERT INTO USERS (username, mail, password) VALUES ($1, $2, $3)",
		username, mail, hashedPassword,
	)
	return err
}

func (r *Repository) FindPasswordByUsername(username string) (string, error) {
	var password string
	err := r.db.QueryRow(
		"SELECT password FROM users WHERE username = $1",
		username,
	).Scan(&password)
	return password, err
}

func (r *Repository) Exists(username string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)",
		username,
	).Scan(&exists)

	return exists, err
}

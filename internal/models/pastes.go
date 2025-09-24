package models

import (
	"database/sql"
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Paste struct {
	ID           int
	UniqueString string
	Content      string
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

type PasteModel struct {
	DB *sql.DB
}

func (m *PasteModel) Insert(uniqueString, content string, passwordHash []byte, expires time.Time) (int, error) {
	query := `INSERT INTO pastes (unique_string, content, password_hash,expires_at)
			  VALUES ($1,$2,$3,$4)
			  RETURNING id`

	var id int

	err := m.DB.QueryRow(query, uniqueString, content, passwordHash, expires).Scan(&id)
	if err != nil {
		return 0, nil
	}
	return id, nil
}

func (m *PasteModel) Get(id int) (*Paste, error) {

	query := `SELECT id, unique_string,content, created_at,expires_at
			  FROM pastes
			  WHERE expires_at > NOW() and id = $1`

	var pastes Paste

	err := m.DB.QueryRow(query, id).Scan(
		&pastes.ID,
		&pastes.UniqueString,
		&pastes.Content,
		&pastes.CreatedAt,
		&pastes.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &pastes, nil
}

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
	PasswordHash []byte
}

type PasteModel struct {
	DB *sql.DB
}

func (m *PasteModel) Insert(uniqueString, content string, passwordHash []byte, expires time.Time) (int, error) {
	query := `INSERT INTO pastes (unique_string, content, password_hash, expires_at)
			  VALUES ($1,$2,$3,$4)
			  RETURNING id`

	var id int
	err := m.DB.QueryRow(query, uniqueString, content, passwordHash, expires).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *PasteModel) Get(id int) (*Paste, error) {
	query := `SELECT id, unique_string, content, created_at, expires_at, password_hash
			  FROM pastes
			  WHERE expires_at > NOW() AND id = $1`

	var paste Paste
	err := m.DB.QueryRow(query, id).Scan(
		&paste.ID,
		&paste.UniqueString,
		&paste.Content,
		&paste.CreatedAt,
		&paste.ExpiresAt,
		&paste.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &paste, nil
}

func (m *PasteModel) GetByUnique(unique string) (*Paste, error) {
	query := `SELECT id, unique_string, content, created_at, expires_at, password_hash
			  FROM pastes
			  WHERE expires_at > NOW() AND unique_string = $1`

	var paste Paste
	err := m.DB.QueryRow(query, unique).Scan(
		&paste.ID,
		&paste.UniqueString,
		&paste.Content,
		&paste.CreatedAt,
		&paste.ExpiresAt,
		&paste.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &paste, nil
}

func (m *PasteModel) Latest(limit int) ([]*Paste, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	query := `SELECT id, unique_string, content, created_at, expires_at, password_hash
			  FROM pastes
			  WHERE expires_at > NOW()
			  ORDER BY created_at DESC
			  LIMIT $1`
	rows, err := m.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	pastes := []*Paste{}
	for rows.Next() {
		var p Paste
		if err := rows.Scan(&p.ID, &p.UniqueString, &p.Content, &p.CreatedAt, &p.ExpiresAt, &p.PasswordHash); err != nil {
			return nil, err
		}
		pastes = append(pastes, &p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return pastes, nil
}

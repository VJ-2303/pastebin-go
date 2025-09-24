CREATE TABLE pastes (
    id BIGSERIAL PRIMARY KEY,
    unique_string TEXT UNIQUE NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    password_hash BYTEA NULL
)

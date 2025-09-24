# pastebin-go

Simple Pastebin-style service written in Go (net/http + julienschmidt/httprouter) with optional password protected pastes and time-limited expiry.

## Features

* Create a paste with content and choose expiry (4h or 24h).
* Optional password. If set, viewers must supply the password to see the content.
* Pastes accessible via nice URL: `/paste/<uniqueID>`.
* Automatic filtering of expired pastes (not returned / 404).

## Schema

```
CREATE TABLE pastes (
	 id BIGSERIAL PRIMARY KEY,
	 unique_string TEXT UNIQUE NOT NULL,
	 content TEXT NOT NULL,
	 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	 expires_at TIMESTAMPTZ NOT NULL,
	 password_hash BYTEA NULL
);
```

## Run Locally

1. Set up Postgres and obtain a DSN, e.g.
	`export DSN="postgres://user:pass@localhost:5432/pastebin?sslmode=disable"`
2. Run migrations (example using psql):
	`psql "$DSN" -f migrations/000001_create_create_pastes_table.up.sql`
3. Start the server:
	`go run ./cmd/web -addr :4000 -dsn "$DSN"`
4. Open `http://localhost:4000/paste/create` to create a paste.

## Password Protected Pastes

If you enter a password while creating a paste, the content will be hashed with bcrypt. When someone visits the paste URL they'll first see a password form. After correct entry, the paste content is shown. Incorrect passwords re-display the form with an error. No sessions/cookies are used; each view requires re-entering the password (easy to extend later with a short-lived cookie).

## Improvements / Next Steps (Not Implemented)

* Add syntax highlighting.
* Add list/recent pastes page (non-passworded only or show lock icon).
* Add rate limiting to creation endpoint.
* Support more flexible expiry options.
* Add tests.

## License

MIT

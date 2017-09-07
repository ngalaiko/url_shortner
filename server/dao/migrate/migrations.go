package migrate

import "time"

type migration struct {
	ID        uint      `db:"id"`
	Name      string    `db:"name"`
	RawSQL    string    `db:"-"`
	FlushSQL  string    `db:"-"`
	AppliedAt time.Time `db:"applied_at"`
}

func initMigrations() []*migration {
	return []*migration{
		{
			Name: "create migrations table",
			RawSQL: `
			CREATE TABLE IF NOT EXISTS migrations (
				id         SERIAL    NOT NULL,
				name       VARCHAR   NOT NULL,
				applied_at TIMESTAMP NOT NULL DEFAULT NOW()
			)
			`,
			FlushSQL: "DELETE FROM migrations",
		},
	}
}

func migrations() []*migration {
	return []*migration{
		{
			Name: "create users table",
			RawSQL: `
			CREATE TABLE users (
				id         BIGSERIAL     NOT NULL PRIMARY KEY,
				first_name VARCHAR(255)  NOT NULL,
				last_name  VARCHAR(255)  NOT NULL,
				created_at TIMESTAMP     NOT NULL,
				deleted_at TIMESTAMP
			)
			`,
			FlushSQL: `DROP TABLE IF EXISTS users`,
		},
		{
			Name: "create links table",
			RawSQL: `
			CREATE TABLE links (
				id         BIGSERIAL NOT NULL PRIMARY KEY,
				user_id    BIGINT    NOT NULL REFERENCES users(id),
				url        TEXT      NOT NULL,
				short_url  TEXT      NOT NULL,
				clicks     BIGINT    NOT NULL DEFAULT 0,
				views      BIGINT    NOT NULL DEFAULT 0,
				expired_at TIMESTAMP NOT NULL,
				created_at TIMESTAMP NOT NULL,
				deleted_at TIMESTAMP
			)
			`,
			FlushSQL: `DROP TABLE IF EXISTS links`,
		},
		{
			Name: "insert user for tests",
			RawSQL: `INSERT INTO users
				(id, first_name, last_name, created_at)
				VALUES
				(0, 'test', 'user', NOW())
			`,
			FlushSQL: "DELETE FROM users WHERE id = 0",
		},
		{
			Name:     "unique links(url, user_id)",
			RawSQL:   `CREATE UNIQUE INDEX links_user_id_url_unique_ix ON links(user_id, url) WHERE deleted_at IS NULL`,
			FlushSQL: `DROP INDEX links_user_id_url_unique_ix`,
		},
		{
			Name:     "unique links(short_url)",
			RawSQL:   `CREATE UNIQUE INDEX links_short_url_unique_ix ON links(short_url) WHERE deleted_at IS NULL`,
			FlushSQL: `DROP INDEX links_short_url_unique_ix`,
		},
	}
}

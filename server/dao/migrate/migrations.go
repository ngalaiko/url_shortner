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
				created_at TIMESTAMP     NOT NULL DEFAULT NOW(),
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
				created_at TIMESTAMP NOT NULL DEFAULT NOW(),
				deleted_at TIMESTAMP
			)
			`,
			FlushSQL: `DROP TABLE IF EXISTS links`,
		},
	}
}

package migrate

import (
	"time"
)

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
					id          BIGSERIAL NOT NULL PRIMARY KEY,
					user_id     BIGINT    NOT NULL REFERENCES users(id),
					url         TEXT      NOT NULL,
					short_url   TEXT      NOT NULL,
					views       BIGINT    NOT NULL DEFAULT 0,
					views_limit BIGINT    NOT NULL DEFAULT 0,
					expired_at  TIMESTAMP NOT NULL,
					created_at  TIMESTAMP NOT NULL,
					deleted_at  TIMESTAMP
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
		{
			Name:     "add facebook id",
			RawSQL:   `ALTER TABLE users ADD COLUMN facebook_id TEXT`,
			FlushSQL: `ALTER TABLE users DROP COLUMN facebook_id`,
		},
		{
			Name:     "add facebook id index",
			RawSQL:   `CREATE UNIQUE INDEX users_facebook_id_unique_ix ON users(facebook_id) WHERE deleted_at IS NULL`,
			FlushSQL: `DROP INDEX users_facebook_id_unique_ix`,
		},
		{
			Name: "tokens table",
			RawSQL: `
				CREATE TABLE user_token (
					id         BIGSERIAL    NOT NULL PRIMARY KEY,
					token      VARCHAR(255) NOT NULL,
					user_id    BIGINT       NOT NULL REFERENCES users(id),
					expired_at TIMESTAMP    NOT NULL
				)
			`,
			FlushSQL: `DROP TABLE user_token`,
		},
		{
			Name: "tokens table unique token index ",
			RawSQL: `
				CREATE UNIQUE INDEX user_token_token_unique_ix ON user_token(token);
			`,
			FlushSQL: `DROP INDEX user_token_token_unique_ix`,
		},
		{
			Name: "drop index links.links_user_id_url_unique_ix",
			RawSQL: `
				DROP INDEX links_user_id_url_unique_ix
			`,
		},
	}
}

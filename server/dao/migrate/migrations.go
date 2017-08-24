package migrate

import "time"

type migration struct {
	ID        uint      `db:"id"`
	Name      string    `db:"name"`
	RawSql    string    `db:"-"`
	FlushSql  string    `db:"-"`
	AppliedAt time.Time `db:"applied_at"`
}

func initMigrations() []*migration {
	return []*migration{
		{
			Name: "create migrations table",
			RawSql: `
			CREATE TABLE IF NOT EXISTS migrations (
				id         SERIAL    NOT NULL,
				name       VARCHAR   NOT NULL,
				applied_at TIMESTAMP NOT NULL DEFAULT NOW()
			)
			`,
		},
	}
}

func migrations() []*migration {
	return []*migration{}
}

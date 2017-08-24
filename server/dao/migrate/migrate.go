package migrate

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/ngalayko/url_shortner/server/dao"
	"github.com/ngalayko/url_shortner/server/logger"
	"go.uber.org/zap"
)

const (
	ctxKey migrationsCtxKey = "migrations_ctx_key"
)

type migrationsCtxKey string

type Migrate struct {
	Db *dao.Db

	logger *logger.Logger
}

func NewContext(ctx context.Context, migrations interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := migrations.(*Migrate); !ok {
		migrations = newMigrate(ctx)
	}

	return context.WithValue(ctx, ctxKey, migrations)
}

func FromContext(ctx context.Context) *Migrate {
	if migrations, ok := ctx.Value(ctxKey).(*Migrate); ok {
		return migrations
	}

	return newMigrate(ctx)
}

func newMigrate(ctx context.Context) *Migrate {
	return &Migrate{
		Db:     dao.FromContext(ctx),
		logger: logger.FromContext(ctx),
	}
}

func (m *Migrate) Apply() error {
	if err := m.applyInitMigrations(); err != nil {
		return err
	}

	appliedMigrations, err := m.applied()
	if err != nil {
		return err
	}

	m.logger.Info("Applying migrations...")

	var count int
	for _, migration := range migrations() {
		if _, ok := appliedMigrations[migration.Name]; ok {
			continue
		}

		if err := m.applyMigration(migration); err != nil {
			return err
		}

		count++
	}

	m.logger.Info("Appled migrations",
		zap.Int("count", count),
	)

	return nil
}

func (m *Migrate) applyInitMigrations() error {
	for _, migration := range initMigrations() {
		if err := m.applyMigration(migration); err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrate) applied() (map[string]*migration, error) {
	var migrations []*migration

	if err := m.Db.Select(&migrations, `
	SELECT m.*
	FROM migrations m
	`); err != nil {
		return nil, err
	}

	result := map[string]*migration{}
	for _, migration := range migrations {
		result[migration.Name] = migration
	}

	return result, nil
}

func (m *Migrate) applyMigration(migration *migration) error {
	return m.Db.Mutate(func(tx *sqlx.Tx) error {

		m.logger.Info("applying migration",
			zap.String("name", migration.Name),
		)

		if _, err := tx.Exec(migration.RawSql); err != nil {
			return err
		}

		if _, err := tx.Exec(`INSERT INTO migrations (name) VALUES ($1)`, migration.Name); err != nil {
			return err
		}

		m.logger.Info("migration applied",
			zap.String("name", migration.Name),
		)

		return nil
	})
}

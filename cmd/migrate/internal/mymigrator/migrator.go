package mymigrator

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/priyanshu-kun/project-metagross/migrations"
)

func NewMigrate(dbURL string) (*migrate.Migrate, error) {
	src, err := iofs.New(migrations.FS, ".")

	if err != nil {
		return nil, fmt.Errorf("load migrations: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, dbURL)

	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return m, nil
}

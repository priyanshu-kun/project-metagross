package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/priyanshu-kun/project-metagross/cmd/migrate/internal/argsparser"
	"github.com/priyanshu-kun/project-metagross/cmd/migrate/internal/commandrunner"
	"github.com/priyanshu-kun/project-metagross/cmd/migrate/internal/migrationsgenerator"
	"github.com/priyanshu-kun/project-metagross/cmd/migrate/internal/mymigrator"
	"github.com/priyanshu-kun/project-metagross/internal/config"
)

const usage = `usage: migrate <command> [args]

commands:
  create <name>     create new up/down migration files
  up [N]            apply all (or N) pending migrations
  down N | --all    revert N migrations (or all)
  force V           set version V and clear dirty flag (no SQL run)
  version           print current version
`

func printVersion(m *migrate.Migrate) {
	v, dirty, err := m.Version()
	switch {
	case errors.Is(err, migrate.ErrNilVersion):
		fmt.Println("Version: no version applied.")
	case err != nil:
		fmt.Fprintln(os.Stderr, "Version:", err)
	default:
		fmt.Printf("version: %d dirty: %v\n", v, dirty)
	}
}

func realMain() int {
	cfg := config.MustLoad()
	argsCfg, err := argsparser.ParseArgs(os.Args[1:], cfg)

	if errors.Is(err, argsparser.ErrHelp) {
		fmt.Print(usage)
		return 0
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		fmt.Fprint(os.Stderr, usage)
		return 2
	}

	switch argsCfg.Cmd {
	case "create", "version", "up", "down", "force":
		// ?????
	default:
		fmt.Fprintf(os.Stderr, "error: Unknown command %q\n", argsCfg.Cmd)
		fmt.Fprint(os.Stderr, usage)
		return 2
	}

	if argsCfg.Cmd == "create" {
		if len(argsCfg.Rest) != 1 {
			fmt.Fprintln(os.Stderr, "error: create needs exactly one name")
			return 2
		}
		if err = migrationsgenerator.Migrationsgenerator(argsCfg.Dir, argsCfg.Rest[0]); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return 1
		}
		return 0
	}

	m, err := mymigrator.NewMigrate(argsCfg.DbURL)

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}
	defer m.Close()

	err = commandrunner.RunCommand(m, argsCfg.Cmd, argsCfg.Rest)

	var short migrate.ErrShortLimit

	switch {
	case err == nil:
	case errors.Is(err, migrate.ErrNoChange):
		fmt.Println("no change")
	case errors.As(err, &short):
		fmt.Println("reached the first migration; nothing more to revert")
	default:
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1

	}

	printVersion(m)

	return 0
}

func main() {
	os.Exit(realMain())
}

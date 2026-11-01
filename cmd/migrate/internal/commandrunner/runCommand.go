package commandrunner

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
)

func RunCommand(m *migrate.Migrate, cmd string, args []string) error {

	switch cmd {
	case "up":
		return m.Up()
	case "down":
		n := 1
		if len(args) > 0 {
			parsed, err := strconv.Atoi(args[0])
			if err != nil || parsed <= 0 {
				return fmt.Errorf("invalid N: %q", args[0])
			}
			n = parsed
		}

		if _, _, err := m.Version(); errors.Is(err, migrate.ErrNilVersion) {
			return migrate.ErrNoChange
		}
		return m.Steps(-n)

	case "force":
		if len(args) == 0 {
			return errors.New("please provide the version name")
		}

		parsed, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid N: %q", args[0])
		}

		return m.Force(parsed)
	case "version":
		return nil
	}

	return fmt.Errorf("unknown command %q", cmd)

}

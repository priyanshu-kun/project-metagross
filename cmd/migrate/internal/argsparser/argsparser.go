package argsparser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/priyanshu-kun/project-metagross/internal/config"
)

type parserConfig struct {
	Dir   string
	DbURL string
	Cmd   string
	Rest  []string
}

var ErrHelp = errors.New("Help requested")

func ParseArgs(args []string, cfg config.Config) (parserConfig, error) {

	parserCfg := parserConfig{
		Dir:   "migrations",
		DbURL: cfg.DB_URL,
	}

	i := 0

	for i < len(args) {

		arg := args[i]

		// fmt.Println(arg)

		if !strings.HasPrefix(arg, "-") {
			break
		}

		name := strings.TrimLeft(arg, "-")

		switch name {
		case "h", "help":
			return parserCfg, ErrHelp
		case "db", "dir":

			if i+1 >= len(args) {
				return parserCfg, fmt.Errorf("flag %s needs a value", arg)
			}
			i++
			if name == "db" {
				parserCfg.DbURL = args[i]
			} else {
				parserCfg.Dir = args[i]
			}

		default:
			return parserCfg, fmt.Errorf("Unknown flag: %s", arg)
		}

		i++
	}

	if i >= len(args) {
		return parserCfg, fmt.Errorf("Missing command")
	}

	parserCfg.Cmd = args[i]
	parserCfg.Rest = args[i+1:]

	return parserCfg, nil
}

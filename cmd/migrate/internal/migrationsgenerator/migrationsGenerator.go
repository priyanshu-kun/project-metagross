package migrationsgenerator

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func Migrationsgenerator(dir string, name string) error {

	version := time.Now().UTC().Format("20060102150405")

	for _, kind := range []string{"up", "down"} {

		path := filepath.Join(dir, fmt.Sprintf("%s_%s.%s.%s", version, name, kind, "sql"))

		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)

		if err != nil {
			return err
		}

		f.Close()

	}

	return nil

}

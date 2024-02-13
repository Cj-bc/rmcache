package main

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Programs map[string]program
}

type program struct {
	Run   []string
	Paths []string
}

// Return true if PROGNAME is found in PATH.
// This requires 'which' command to run.
func which(progname string) bool {
	err := exec.Command("which", progname).Run()
	return err == nil
}

func main() {
	logger := slog.Default()

	var cfg *Config
	_, err := toml.DecodeFile("example.toml", &cfg)
	if err != nil {
		logger.Error("Failed to decode file", "error", err)
		os.Exit(-1)
	}

	for name, program := range cfg.Programs {
		if !which(name) {
			logger.Info("Program isn't found. Skipping entry", "program", name)
			continue
		}
		logger.Info("Found entry", "program", name)

		if len(program.Run) != 0 {
			// Run it
		}

		for _, glob := range program.Paths {
			expanded, err := filepath.Glob(glob)
			if err != nil {
				logger.Info("Glob pattern malformed.", "path", glob, "err", err)
				continue
			}

			for _, p := range expanded {
				l := logger.With("path", p, "program", name)
				err := removePath(p)
				if err != nil {
					l.Warn("Failed to remove path", "error", err)
				} else {
					l.Info("Proprely removed path")
				}
			}
		}

	}
}

// Remove given path if it is regular file. This doesn't allow to
// specify symbolic link for security reason
func removePath(path string) error {
	if info, err := os.Stat(path); err != nil {
		return err
	} else if !info.Mode().IsRegular() {
		return errors.New("File is not regular file")
	}

	return os.Remove(path)
}

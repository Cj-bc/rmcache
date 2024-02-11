package main

import (
	"log/slog"
	"os"
	"os/exec"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Programs map[string]program
}

type program struct {
	Run   []string
	Paths []string
}

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

		if len(program.Run) != 0 {
			// Run it
		}

		for _, p := range program.Paths {
			logger.Info(p)
			err := removePath(p)
			if err != nil {
				logger.Warn("Failed to remove path", "path", p)
			} else {
				logger.Info("Proprely removed path", "program", name, "path", p)
			}
		}

	}
}

// Remove given path by `unlink'
// TODO:
// - make sure 'p' is absolute path
// - return error if 'p' is symlink (for security reason)
func removePath(p string) error {
	err := unix.Unlink(p)
	return err
}

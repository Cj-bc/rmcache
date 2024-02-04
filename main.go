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

	}
}

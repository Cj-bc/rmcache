package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
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
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		logger.Warn("Could not locate user config dir.")
	}
	var (
		config string
		dryRun bool
	)

	app := &cli.App{
		Name: "rmcache",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:        "config",
				Value:       filepath.Join(userConfigDir, "rmcache", "config.toml"),
				Destination: &config,
			},
			&cli.BoolFlag{
				Name:        "dryrun",
				Aliases:     []string{"d"},
				Destination: &dryRun,
			},
		},
		Action: func(cCtx *cli.Context) error {
			var cfg *Config
			_, err = toml.DecodeFile(config, &cfg)
			if err != nil {
				logger.Error("Failed to decode file", "error", err)
				return err
			}

			for name, program := range cfg.Programs {
				if !which(name) {
					logger.Info("Program isn't found. Skipping entry", "program", name)
					continue
				}
				logger.Info("Found entry", "program", name)

				if len(program.Run) != 0 {
					logger.Info("About to run program", "executable", program.Run[0], "args", strings.Join(program.Run[1:], " "))
					// Run it
				}

				for _, glob := range program.Paths {
					expanded, err := Expands(glob)
					if err != nil {
						logger.Info("Error while expanding glob pattern. Skipping path", "path", glob, "error", err)
						continue
					}

					for _, p := range expanded {
						l := logger.With("path", p, "program", name)

						l.Info("Trying to remove file")

						info, err := os.Stat(p)
						if err != nil {
							l.Warn("Cannot stat(2) file", "error", err)
						} else if !info.Mode().IsRegular() {
							l.Warn("File is not regular file", "error", err)
						}

						if !dryRun {
							err = os.Remove(p)
						} else {
							err = nil
						}
						if err != nil {
							l.Warn("Failed to remove path", "error", err)
						} else {
							l.Info("Proprely removed path", "size(bytes)", info.Size())
						}
					}
				}
			}

			return nil
		},
	}

	app.Run(os.Args)

}

// Expands tildas, globs
func Expands(path string) ([]string, error) {
	_path := path
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return []string{}, err
		}

		_path = filepath.Join(home, path[2:])
	}

	expanded, err := filepath.Glob(_path)
	if err != nil {
		return []string{}, fmt.Errorf("malformed pattern: %e", err)
	} else if expanded == nil {
		return []string{}, errors.New("Glob pattern matches nothing.")
	}

	return expanded, nil
}

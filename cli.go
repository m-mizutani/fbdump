package main

import (
	"github.com/m-mizutani/zlog"
	"github.com/urfave/cli/v2"
)

type config struct {
	RepoDir string
}

var logger = zlog.New()

func Run(argv []string) error {
	var logLevel string
	cfg := &config{}

	app := &cli.App{
		Name:  "fbdump",
		Usage: "Dump Firebase Auth users",
		Commands: []*cli.Command{
			cmdDump(cfg),
			cmdLoad(cfg),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "repo-dir",
				Usage:       "Repository directory",
				Value:       "repo",
				Aliases:     []string{"r"},
				EnvVars:     []string{"FBDUMP_REPO_DIR"},
				Destination: &cfg.RepoDir,
			},
			&cli.StringFlag{
				Name:        "log-level",
				Aliases:     []string{"l"},
				Usage:       "Log level [debug|info|warn|error]",
				EnvVars:     []string{"FBDUMP_LOG_LEVEL"},
				Destination: &logLevel,
				Value:       "info",
			},
		},
		Before: func(ctx *cli.Context) error {
			logger = logger.Clone(zlog.WithLogLevel(logLevel))
			return nil
		},
	}

	if err := app.Run(argv); err != nil {
		logger.Err(err).Error("exit with error")
		return err
	}

	return nil
}

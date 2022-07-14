package main

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v2"
)

func cmdLoad(cfg *config) *cli.Command {
	var output string
	return &cli.Command{
		Name:    "load",
		Aliases: []string{"l"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "Output file, '-' means stdout",
				Value:       "-",
				EnvVars:     []string{"FBDUMP_OUT"},
				Destination: &output,
			},
		},
		Action: func(ctx *cli.Context) error {
			repo := NewRepository(cfg.RepoDir, 0)

			out := os.Stdout
			if output != "-" {
				fd, err := os.Create(output)
				if err != nil {
					return goerr.Wrap(err, "create output file")
				}
				defer fd.Close()
				out = fd
			}

			if err := load(ctx.Context, repo, out); err != nil {
				return err
			}

			return nil
		},
	}
}

func load(ctx context.Context, repo *Repository, out io.WriteCloser) error {
	encoder := json.NewEncoder(out)

	recCh, errCh := repo.Load()
	for {
		select {
		case record := <-recCh:
			if record == nil {
				return nil
			}
			if err := encoder.Encode(record); err != nil {
				return goerr.Wrap(err, "encode record")
			}

		case err := <-errCh:
			return err
		}
	}
}

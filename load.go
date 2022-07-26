package main

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"firebase.google.com/go/v4/auth"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/gots/slice"
	"github.com/urfave/cli/v2"
)

func cmdLoad(cfg *config) *cli.Command {
	var (
		output    string
		providers cli.StringSlice
	)
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
			&cli.StringSliceFlag{
				Name:        "provider",
				Aliases:     []string{"p"},
				Usage:       "Specify provider(s) to filter output [password|google.com|apple.com|...]",
				Destination: &providers,
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

			if err := load(ctx.Context, repo, providers.Value(), out); err != nil {
				return err
			}

			return nil
		},
	}
}

func matchProvider(providers []string, userInfo []*auth.UserInfo) bool {
	if len(providers) == 0 {
		return true // path through if providers is not defined
	}

	for _, info := range userInfo {
		if slice.Contains(providers, info.ProviderID) {
			return true
		}
	}
	return false
}

func load(ctx context.Context, repo *Repository, providers []string, out io.WriteCloser) error {
	encoder := json.NewEncoder(out)

	recCh, errCh := repo.Load()
	for {
		select {
		case record := <-recCh:
			if record == nil {
				return nil
			}

			if !matchProvider(providers, record.ProviderUserInfo) {
				break
			}

			if err := encoder.Encode(record); err != nil {
				return goerr.Wrap(err, "encode record")
			}

		case err := <-errCh:
			return err
		}
	}
}

package main

import (
	"context"
	"path/filepath"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/m-mizutani/goerr"
	"github.com/urfave/cli/v2"
	"google.golang.org/api/iterator"
)

func cmdDump(cfg *config) *cli.Command {
	var (
		projectID     string
		stateFilePath string
		depth         int
	)

	return &cli.Command{
		Name:    "dump",
		Aliases: []string{"d"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "project-id",
				Usage:       "Google Cloud Project ID",
				EnvVars:     []string{"FBDUMP_PROJECT_ID"},
				Destination: &projectID,
			},
			&cli.StringFlag{
				Name:        "state",
				Usage:       "State file path",
				EnvVars:     []string{"FBDUMP_STATE"},
				Value:       "state.json",
				Destination: &stateFilePath,
			},
			&cli.IntFlag{
				Name:        "depth",
				Usage:       "depth of dump directory structure",
				EnvVars:     []string{"FBDUMP_DEPTH"},
				Value:       2,
				Destination: &depth,
			},
		},
		Action: func(ctx *cli.Context) error {
			state, err := NewState(filepath.Clean(stateFilePath))
			if err != nil {
				return err
			}

			repo := NewRepository(cfg.WorkDir, depth)

			if err := dump(ctx.Context, state, repo, projectID); err != nil {
				return err
			}
			return nil
		},
	}
}

const firebaseListUsersLimit = 1000

func dump(ctx context.Context, state *State, repo *Repository, projectID string) error {
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: projectID,
	})
	if err != nil {
		return goerr.Wrap(err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return goerr.Wrap(err)
	}

	iter := client.Users(ctx, "")
	pager := iterator.NewPager(iter, firebaseListUsersLimit, string(state.PageToken))

	for {
		var users []*auth.ExportedUserRecord
		nextToken, err := pager.NextPage(&users)
		if err != nil {
			return goerr.Wrap(err)
		}
		if len(users) == 0 {
			break
		}

		if err := repo.Put(state.PageToken, users); err != nil {
			return err
		}

		logger.With("size", len(users)).With("current", state.PageToken).With("next", nextToken).Debug("Read a page")

		if err := state.Update(PageToken(nextToken)); err != nil {
			return err
		}
	}

	return nil
}

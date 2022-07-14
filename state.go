package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/m-mizutani/goerr"
)

type State struct {
	PageToken     PageToken `json:"page_token"`
	LastUpdatedAt time.Time `json:"last_updated_at"`

	path string
}

func NewState(filePath string) (*State, error) {
	path := filepath.Clean(filePath)

	f, err := os.Open(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, goerr.Wrap(err)
		}

		f, err := os.Create(path)
		if err != nil {
			return nil, goerr.Wrap(err)
		}
		defer f.Close()

		state := &State{
			path:          filePath,
			LastUpdatedAt: time.Now(),
		}
		if err := json.NewEncoder(f).Encode(&state); err != nil {
			return nil, goerr.Wrap(err, "failed to create state")
		}

		return state, nil
	}

	var state State
	if err := json.NewDecoder(f).Decode(&state); err != nil {
		return nil, goerr.Wrap(err, "failed to read state")
	}
	state.path = path

	return &state, nil
}

func (x *State) Update(token PageToken) error {
	f, err := os.Create(x.path)
	if err != nil {
		return goerr.Wrap(err)
	}
	defer f.Close()

	x.PageToken = token
	x.LastUpdatedAt = time.Now()

	if err := json.NewEncoder(f).Encode(x); err != nil {
		return goerr.Wrap(err, "failed to update state")
	}

	return nil
}

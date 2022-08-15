package main

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"firebase.google.com/go/v4/auth"
	"github.com/m-mizutani/goerr"
)

type Repository struct {
	depth   int
	baseDir string
}

func NewRepository(baseDir string, depth int) *Repository {
	if depth < 0 {
		panic("depth must be equal or greater than 0")
	}

	return &Repository{
		depth:   depth,
		baseDir: filepath.Clean(baseDir),
	}
}

func (x *Repository) Put(token PageToken, users []*auth.ExportedUserRecord) error {
	var fpath string
	if len(users) == 0 {
		return nil
	}
	if len(token) == 0 { // use UID of head of users if token is empty (first request)
		if err := os.MkdirAll(x.baseDir, 0755); err != nil {
			return goerr.Wrap(err)
		}
		fpath = filepath.Join(x.baseDir, "head.json.gz")
	} else {
		dirs := []string{x.baseDir}
		for i := 0; i < len(token) && i < x.depth; i++ {
			dirs = append(dirs, string(token[i]))
		}

		dir := filepath.Join(dirs...)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return goerr.Wrap(err)
		}
		fname := token + ".json.gz"
		fpath = filepath.Join(dir, filepath.Clean(string(fname)))
	}

	fd, err := os.Create(fpath)
	if err != nil {
		return goerr.Wrap(err)
	}
	defer func() {
		if err := fd.Close(); err != nil {
			logger.Err(err).Warn("failed to close file writer")
		}
	}()

	w := gzip.NewWriter(fd)
	defer func() {
		if err := w.Close(); err != nil {
			logger.Err(err).Warn("failed to close gzip writer")
		}
	}()

	if err := json.NewEncoder(w).Encode(users); err != nil {
		return goerr.Wrap(err, "encode exported users to json")
	}
	logger.With("path", fpath).Debug("saved user records")

	return nil
}

func (x *Repository) Load() (chan *auth.ExportedUserRecord, chan error) {
	recCh := make(chan *auth.ExportedUserRecord)
	errCh := make(chan error)

	go func() {
		defer close(recCh)
		defer close(errCh)

		if err := filepath.WalkDir(x.baseDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return goerr.Wrap(err).With("path", path)
			}
			if d.IsDir() {
				return nil
			}

			users := []*auth.ExportedUserRecord{}
			fd, err := os.Open(path)
			if err != nil {
				return goerr.Wrap(err, "file open").With("path", path)
			}
			defer fd.Close()

			var rd io.Reader = fd
			if filepath.Ext(path) == ".gz" {
				gz, err := gzip.NewReader(fd)
				if err != nil {
					return goerr.Wrap(err, "new gzip reader")
				}
				rd = gz
			}

			if err := json.NewDecoder(rd).Decode(&users); err != nil {
				return goerr.Wrap(err)
			}

			for i := range users {
				recCh <- users[i]
			}

			return nil
		}); err != nil {
			errCh <- err
			return
		}
	}()

	return recCh, errCh
}

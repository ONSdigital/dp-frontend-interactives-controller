package storage

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"io"
	"net/http"
	"os"
)

func NewLocalFilesystemProvider(root http.Dir) localfs {
	return localfs{root: root}
}

type localfs struct {
	root http.Dir
}

func (s localfs) Get(path string) (io.ReadCloser, error) {
	file, err := os.Open(fmt.Sprintf("%s%s", s.root, path))
	if err != nil {
		return nil, fmt.Errorf("cannot open file %w", err)
	}

	return file, err
}

func (s localfs) Checker() func(context.Context, *healthcheck.CheckState) error {
	return func(_ context.Context, s *healthcheck.CheckState) error {
		return s.Update(healthcheck.StatusOK, "localfs healthy", http.StatusOK)
	}
}

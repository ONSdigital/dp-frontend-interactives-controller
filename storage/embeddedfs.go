package storage

import (
	"context"
	"embed"
	"fmt"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"io"
	"net/http"
)

var (
	//go:embed embeddedfs/*
	static embed.FS
)

func NewFromEmbeddedFilesystem() embeddedfs {
	return embeddedfs{}
}

type embeddedfs struct{}

func (s embeddedfs) Get(path string) (io.ReadCloser, error) {
	file, err := static.Open(fmt.Sprintf("%s%s", "embeddedfs", path))
	if err != nil {
		return nil, fmt.Errorf("cannot open file %w", err)
	}

	return file, err
}

func (s embeddedfs) Checker() func(context.Context, *healthcheck.CheckState) error {
	return func(_ context.Context, s *healthcheck.CheckState) error {
		return s.Update(healthcheck.StatusOK, "embeddedfs healthy", http.StatusOK)
	}
}

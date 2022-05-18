package storage

import (
	"context"
	"io"
	"net/http"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

func NewFromDownloadService(client DownloadServiceAPIClient) downloadService {
	return downloadService{client}
}

type downloadService struct {
	client DownloadServiceAPIClient
}

func (s downloadService) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	resp, err := s.client.Download(ctx, path)
	if err != nil {
		return nil, err
	}
	return resp.Content, nil
}

func (s downloadService) Checker() func(context.Context, *healthcheck.CheckState) error {
	return func(_ context.Context, s *healthcheck.CheckState) error {
		//todo sort this when ONS client done
		return s.Update(healthcheck.StatusOK, "download-service healthy", http.StatusOK)
	}
}

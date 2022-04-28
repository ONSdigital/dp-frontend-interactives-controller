package storage

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"io"
	"net/http"
)

func NewFromDownloadService(serviceAuthToken string, client DownloadServiceAPIClient) downloadService {
	return downloadService{serviceAuthToken, client}
}

type downloadService struct {
	serviceAuthToken string
	client           DownloadServiceAPIClient
}

func (s downloadService) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	resp, err := s.client.Download(ctx, "", s.serviceAuthToken, path)
	if err != nil {
		return nil, err
	}
	if resp.RedirectUrl != "" {
		return nil, fmt.Errorf("%s redirecting to %s", path, resp.RedirectUrl)
	}
	return resp.Content, nil
}

func (s downloadService) Checker() func(context.Context, *healthcheck.CheckState) error {
	return func(_ context.Context, s *healthcheck.CheckState) error {
		//todo sort this when ONS client done
		return s.Update(healthcheck.StatusOK, "download-service healthy", http.StatusOK)
	}
}

package storage

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/download"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"io"
	"net/http"
)

func NewFromDownloadService(serviceAuthToken string, client *download.Client) downloadService {
	return downloadService{serviceAuthToken, client}
}

type downloadService struct {
	serviceAuthToken string
	client           *download.Client
}

func (s downloadService) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	return s.client.Download(ctx, "", s.serviceAuthToken, path)
}

func (s downloadService) Checker() func(context.Context, *healthcheck.CheckState) error {
	return func(_ context.Context, s *healthcheck.CheckState) error {
		//todo sort this when ONS client done
		return s.Update(healthcheck.StatusOK, "download-service healthy", http.StatusOK)
	}
}

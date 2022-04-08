package storage

import (
	"context"
	openapi "github.com/ONSdigital/dp-frontend-interactives-controller/internal/dp-download-service-client/go"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"io"
	"net/http"
)

func NewFromDownloadService(serviceAuthToken, host, scheme string) downloadService {
	cfg := openapi.NewConfiguration()
	cfg.Host = host
	cfg.Scheme = scheme
	cfg.AddDefaultHeader("Authorization", "Bearer "+serviceAuthToken)
	return downloadService{client: openapi.NewAPIClient(cfg)}
}

type downloadService struct {
	client *openapi.APIClient
}

func (s downloadService) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	req := s.client.DownloadFileApi.DownloadsNewFilepathGet(ctx, path)
	_, resp, err := s.client.DownloadFileApi.DownloadsNewFilepathGetExecute(req)
	if resp.StatusCode == 200 {
		//ignore undefined response type errors
		return resp.Body, nil
	}
	return resp.Body, err
}

func (s downloadService) Checker() func(context.Context, *healthcheck.CheckState) error {
	return func(_ context.Context, s *healthcheck.CheckState) error {
		//todo sort this when ONS client done
		return s.Update(healthcheck.StatusOK, "download-service healthy", http.StatusOK)
	}
}

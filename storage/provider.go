package storage

import (
	"context"
	"io"

	"github.com/ONSdigital/dp-api-clients-go/v2/download"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

//go:generate moq -out mocks/provider.go -pkg mocks_storage . Provider S3Bucket DownloadServiceAPIClient

type Provider interface {
	Get(context.Context, string) (io.ReadCloser, error)
	Checker() func(context.Context, *healthcheck.CheckState) error
}

// S3Bucket defines methods used from dp-s3 lib - init points to a specific bucket
type S3Bucket interface {
	Get(key string) (io.ReadCloser, *int64, error)
	Checker(ctx context.Context, state *healthcheck.CheckState) (err error)
}

type DownloadServiceAPIClient interface {
	Download(ctx context.Context, path string) (*download.Response, error)
	Checker(ctx context.Context, state *healthcheck.CheckState) (err error)
}

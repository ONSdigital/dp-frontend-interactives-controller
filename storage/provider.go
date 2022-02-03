package storage

import (
	"context"
	"io"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

//go:generate moq -out mocks/provider.go -pkg mocks_storage . Provider
//go:generate moq -out mocks/s3bucket.go -pkg mocks_storage . S3Bucket

type Provider interface {
	Get(path string) (io.ReadCloser, error)
	Checker() func(context.Context, *healthcheck.CheckState) error
}

// S3Bucket defines methods used from dp-s3 lib - init points to a specific bucket
type S3Bucket interface {
	Get(key string) (io.ReadCloser, *int64, error)
	Checker(ctx context.Context, state *healthcheck.CheckState) (err error)
}

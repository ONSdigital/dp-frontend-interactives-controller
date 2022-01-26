package storage

import (
	"context"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"io"
)

//go:generate moq -out mocks/s3bucket.go -pkg mocks_storage . S3Bucket

// S3Bucket defines methods used from dp-s3 lib - init points to a specific bucket
type S3Bucket interface {
	Get(key string) (io.ReadCloser, *int64, error)
	Checker(ctx context.Context, state *healthcheck.CheckState) (err error)
}

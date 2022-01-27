package storage

import (
	"context"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"io"
)

func NewS3Provider(s3bucket S3Bucket) s3 {
	return s3{s3bucket: s3bucket}
}

type s3 struct {
	s3bucket S3Bucket
}

func (s s3) Get(path string) (io.ReadCloser, error) {
	s3ReadCloser, _, err := s.s3bucket.Get(path)
	return s3ReadCloser, err
}

func (s s3) Checker() func(context.Context, *healthcheck.CheckState) error {
	return s.s3bucket.Checker
}

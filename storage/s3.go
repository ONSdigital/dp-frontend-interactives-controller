package storage

import (
	"context"
	"io"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

func NewFromS3Bucket(s3bucket S3Bucket) s3 {
	return s3{s3bucket: s3bucket}
}

type s3 struct {
	s3bucket S3Bucket
}

func (s s3) Get(_ context.Context, path string) (io.ReadCloser, error) {
	s3ReadCloser, _, err := s.s3bucket.Get(path)
	return s3ReadCloser, err
}

func (s s3) Checker() func(context.Context, *healthcheck.CheckState) error {
	return s.s3bucket.Checker
}

//func (e *Init) DoGetS3Bucket() (storage.S3Bucket, error) {
//	//TODO awsify this...e.g. https://github.com/ONSdigital/dp-download-service/blob/c750eae9e7eaea003420aec432bc9a7322a3782c/service/external/external.go#L49
//	//if cfg.LocalObjectStore != "" {
//	s3Config := &aws.Config{
//		Credentials:      credentials.NewStaticCredentials("minio-access-key", "minio-secret-key", ""),
//		Endpoint:         aws.String("http://localhost:9001"),
//		Region:           aws.String("eu-west-1"),
//		DisableSSL:       aws.Bool(true),
//		S3ForcePathStyle: aws.Bool(true),
//	}
//
//	s, err := session.NewSession(s3Config)
//	if err != nil {
//		return nil, fmt.Errorf("could not create s3 session: %w", err)
//	}
//	s3 := s3client.NewClientWithSession("private-bucket", s)
//	//}
//
//	//s3, err := s3client.NewClient(cfg.AwsRegion, cfg.BucketName)
//	//if err != nil {
//	//	return nil, fmt.Errorf("could not create s3 client: %w", err)
//	//}
//
//	return s3, nil
//}

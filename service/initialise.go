package service

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-frontend-interactives-controller/config"
	"github.com/ONSdigital/dp-frontend-interactives-controller/storage"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/http"
	s3client "github.com/ONSdigital/dp-s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"net/http"
)

// ExternalServiceList holds the initialiser and initialisation state of external services.
type ExternalServiceList struct {
	HealthCheck bool
	Init        Initialiser
}

// NewServiceList creates a new service list with the provided initialiser
func NewServiceList(initialiser Initialiser) *ExternalServiceList {
	return &ExternalServiceList{
		Init: initialiser,
	}
}

// Init implements the Initialiser interface to initialise dependencies
type Init struct{}

// GetHTTPServer creates an http server
func (e *ExternalServiceList) GetHTTPServer(bindAddr string, router http.Handler) HTTPServer {
	s := e.Init.DoGetHTTPServer(bindAddr, router)
	return s
}

// GetHandlers creates handlers depending on config: localfs, s3, (todo) static-file-service
func (e *ExternalServiceList) GetStorageProvider(cfg *config.Config) (storage.Provider, error) {
	var sp storage.Provider

	if len(cfg.ServeFromLocalDir) > 0 {
		sp = storage.NewLocalFilesystemProvider(http.Dir(cfg.ServeFromLocalDir))
	} else {
		sourceS3bucket, err := e.Init.DoGetS3Bucket()
		if err != nil {
			return nil, fmt.Errorf("could not get s3 bucket: %w", err)
		}

		sp = storage.NewS3Provider(sourceS3bucket)
	}

	return sp, nil
}

// GetHealthClient returns a healthclient for the provided URL
func (e *ExternalServiceList) GetHealthClient(name, url string) *health.Client {
	return e.Init.DoGetHealthClient(name, url)
}

// GetHealthCheck creates a healthcheck with versionInfo and sets the HealthCheck flag to true
func (e *ExternalServiceList) GetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	hc, err := e.Init.DoGetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		return nil, err
	}
	e.HealthCheck = true
	return hc, nil
}

// DoGetHTTPServer creates an HTTP Server with the provided bind address and router
func (e *Init) DoGetHTTPServer(bindAddr string, router http.Handler) HTTPServer {
	s := dphttp.NewServer(bindAddr, router)
	s.HandleOSSignals = false
	return s
}

// DoGetS3Bucket obtains a new S3 bucket client, or minio client if a non-empty LocalObjectStore is provided
func (e *Init) DoGetS3Bucket() (storage.S3Bucket, error) {
	//TODO awsify this...e.g. https://github.com/ONSdigital/dp-download-service/blob/c750eae9e7eaea003420aec432bc9a7322a3782c/service/external/external.go#L49
	//if cfg.LocalObjectStore != "" {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials("minio-access-key", "minio-secret-key", ""),
		Endpoint:         aws.String("http://localhost:9001"),
		Region:           aws.String("eu-west-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	s, err := session.NewSession(s3Config)
	if err != nil {
		return nil, fmt.Errorf("could not create s3 session: %w", err)
	}
	s3 := s3client.NewClientWithSession("private-bucket", s)
	//}

	//s3, err := s3client.NewClient(cfg.AwsRegion, cfg.BucketName)
	//if err != nil {
	//	return nil, fmt.Errorf("could not create s3 client: %w", err)
	//}

	return s3, nil
}

// DoGetHealthClient creates a new Health Client for the provided name and url
func (e *Init) DoGetHealthClient(name, url string) *health.Client {
	return health.NewClient(name, url)
}

// DoGetHealthCheck creates a healthcheck with versionInfo
func (e *Init) DoGetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	versionInfo, err := healthcheck.NewVersionInfo(buildTime, gitCommit, version)
	if err != nil {
		return nil, err
	}
	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	return &hc, nil
}

// NewMockHTTPClient mocks HTTP Client
func NewMockHTTPClient(r *http.Response, err error) *dphttp.ClienterMock {
	return &dphttp.ClienterMock{
		SetPathsWithNoRetriesFunc: func(paths []string) {},
		GetPathsWithNoRetriesFunc: func() []string { return []string{} },
		DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return r, err
		},
	}
}

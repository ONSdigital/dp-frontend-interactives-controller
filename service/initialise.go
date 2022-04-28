package service

import (
	"github.com/ONSdigital/dp-api-clients-go/v2/download"
	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-api-clients-go/v2/interactives"
	"github.com/ONSdigital/dp-frontend-interactives-controller/config"
	"github.com/ONSdigital/dp-frontend-interactives-controller/routes"
	"github.com/ONSdigital/dp-frontend-interactives-controller/storage"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/http"
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

func (e *ExternalServiceList) GetDownloadServiceAPIClient(cfg *config.Config) (storage.DownloadServiceAPIClient, error) {
	if cfg.ServeFromEmbeddedContent {
		return nil, nil
	}
	return e.Init.DoGetDownloadServiceAPIClient(cfg)
}

// GetStorageProvider returns storage provider depending on config: localfs, s3, static files (dp-download-service)
func (e *ExternalServiceList) GetStorageProvider(cfg *config.Config, c storage.DownloadServiceAPIClient) (storage.Provider, error) {
	return e.Init.DoGetStorageProvider(cfg, c)
}

// GetInteractivesAPIClient creates an interactives api client and sets the InteractivesApi flag to true
func (e *ExternalServiceList) GetInteractivesAPIClient(apiRouter *health.Client) (routes.InteractivesAPIClient, error) {
	client, err := e.Init.DoGetInteractivesAPIClient(apiRouter)
	if err != nil {
		return nil, err
	}

	return client, nil
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

// DoGetInteractivesApiClient returns an interactives api client
func (e *Init) DoGetInteractivesAPIClient(apiRouter *health.Client) (routes.InteractivesAPIClient, error) {
	apiClient := interactives.NewWithHealthClient(apiRouter, "v1")
	return apiClient, nil
}

func (e *Init) DoGetDownloadServiceAPIClient(cfg *config.Config) (storage.DownloadServiceAPIClient, error) {
	apiClient := download.NewAPIClient(cfg.DownloadAPIURL)
	return apiClient, nil
}

func (e *Init) DoGetStorageProvider(cfg *config.Config, downloadClient storage.DownloadServiceAPIClient) (storage.Provider, error) {
	var sp storage.Provider

	if cfg.ServeFromEmbeddedContent {
		sp = storage.NewFromEmbeddedFilesystem()
	} else {
		sp = storage.NewFromDownloadService(cfg.ServiceAuthToken, downloadClient)
	}

	return sp, nil
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

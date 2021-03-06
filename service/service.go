package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-frontend-interactives-controller/config"
	"github.com/ONSdigital/dp-frontend-interactives-controller/handlers"
	"github.com/ONSdigital/dp-frontend-interactives-controller/routes"
	"github.com/ONSdigital/dp-frontend-interactives-controller/storage"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

// Service contains the healthcheck, server and serviceList for the controller
type Service struct {
	Config                   *config.Config
	HealthCheck              HealthChecker
	Server                   HTTPServer
	ServiceList              *ExternalServiceList
	StorageProvider          storage.Provider
	APIRouter                *health.Client
	DownloadServiceAPIClient storage.DownloadServiceAPIClient
}

// New creates a new service
func New(cfg *config.Config, serviceList *ExternalServiceList) *Service {
	return &Service{
		Config:      cfg,
		ServiceList: serviceList,
	}
}

// Init initialises all the service dependencies, including healthcheck with checkers, api and middleware
func (s *Service) Init(ctx context.Context) (err error) {

	// API router & clients
	s.APIRouter = s.ServiceList.GetHealthClient("api-router", s.Config.APIRouterURL)
	interactivesAPIClient, err := s.ServiceList.GetInteractivesAPIClient(s.APIRouter)
	if err != nil {
		return err
	}

	s.DownloadServiceAPIClient, err = s.ServiceList.GetDownloadServiceAPIClient(s.Config)
	if err != nil {
		return err
	}

	// Init storage provider
	s.StorageProvider, err = s.ServiceList.GetStorageProvider(s.Config, s.DownloadServiceAPIClient)
	if err != nil {
		return fmt.Errorf("failed to initialise storage provider %w", err)
	}

	// Init healthcheck with checkers for downstream deps (do this after initing any deps that need checking!)
	s.HealthCheck, err = s.ServiceList.GetHealthCheck(s.Config, BuildTime, GitCommit, Version)
	if err != nil {
		return fmt.Errorf("failed to create health check %w", err)
	}
	if err = s.registerCheckers(ctx); err != nil {
		return fmt.Errorf("failed to register checkers %w", err)
	}

	// Init clients
	clients := routes.Clients{
		Storage: s.StorageProvider,
		API:     interactivesAPIClient,
	}

	// Init router
	r := mux.NewRouter()
	routes.Setup(s.Config, r, s.HealthCheck.Handler,
		handlers.Interactives(s.Config, clients),
		handlers.InteractivesRedirect(s.Config, clients))
	s.Server = s.ServiceList.GetHTTPServer(s.Config.BindAddr, r)

	return nil
}

// Run starts an initialised service
func (s *Service) Run(ctx context.Context, svcErrors chan error) {
	log.Info(ctx, "Starting service", log.Data{"config": s.Config})

	// Start healthcheck
	s.HealthCheck.Start(ctx)

	// Start HTTP server
	log.Info(ctx, "Starting server")
	go func() {
		if err := s.Server.ListenAndServe(); err != nil {
			svcErrors <- fmt.Errorf("failed to start http listen and serve %w", err)
		}
	}()
}

// Close gracefully shuts the service down in the required order, with timeout
func (s *Service) Close(ctx context.Context) error {
	log.Info(ctx, "commencing graceful shutdown")
	ctx, cancel := context.WithTimeout(ctx, s.Config.GracefulShutdownTimeout)
	hasShutdownError := false

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		log.Info(ctx, "stop health checkers")
		s.HealthCheck.Stop()

		// TODO: close any backing services here, e.g. client connections to databases

		// stop any incoming requests
		if err := s.Server.Shutdown(ctx); err != nil {
			log.Error(ctx, "failed to shutdown http server", err)
			hasShutdownError = true
		}
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		log.Error(ctx, "shutdown timed out", ctx.Err())
		return ctx.Err()
	}

	// other error
	if hasShutdownError {
		err := errors.New("failed to shutdown gracefully")
		log.Error(ctx, "failed to shutdown gracefully ", err)
		return err
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func (s *Service) registerCheckers(ctx context.Context) (err error) {
	hasErrors := false

	if err = s.HealthCheck.AddCheck("storage provider", s.StorageProvider.Checker()); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for storage provider", err)
	}

	if err = s.HealthCheck.AddCheck("API router", s.APIRouter.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "failed to add API router health checker", err)
	}

	if s.DownloadServiceAPIClient != nil {
		if err = s.HealthCheck.AddCheck("DownloadService API", s.DownloadServiceAPIClient.Checker); err != nil {
			hasErrors = true
			log.Error(ctx, "failed to add DownloadService API health checker", err)
		}
	}

	if hasErrors {
		return errors.New("error(s) registering checkers for healthcheck")
	}

	return nil
}

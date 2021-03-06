package service_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ONSdigital/dp-api-clients-go/v2/download"
	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-api-clients-go/v2/interactives"
	"github.com/ONSdigital/dp-frontend-interactives-controller/config"
	"github.com/ONSdigital/dp-frontend-interactives-controller/routes"
	mocks_routes "github.com/ONSdigital/dp-frontend-interactives-controller/routes/mocks"
	"github.com/ONSdigital/dp-frontend-interactives-controller/service"
	mocks_service "github.com/ONSdigital/dp-frontend-interactives-controller/service/mocks"
	"github.com/ONSdigital/dp-frontend-interactives-controller/storage"
	mocks_storage "github.com/ONSdigital/dp-frontend-interactives-controller/storage/mocks"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/http"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	expectedCheckCalls = 3
)

var (
	ctx = context.Background()

	errAddCheckFail = errors.New("error(s) registering checkers for healthcheck")
	errHealthCheck  = errors.New("healthCheck error")
	errServer       = errors.New("HTTP Server error")

	// Health Check Mock
	hcMock = &mocks_service.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
	}
	hcMockAddFail = &mocks_service.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return errAddCheckFail },
		StartFunc:    func(ctx context.Context) {},
	}
	funcDoGetHealthCheckOK = func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
		return hcMock, nil
	}
	funcDoGetHealthCheckFail = func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
		return nil, errHealthCheck
	}
	funcDoGetHealthAddCheckerFail = func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
		return hcMockAddFail, nil
	}

	// Server Mock
	serverWg   = &sync.WaitGroup{}
	serverMock = &mocks_service.HTTPServerMock{
		ListenAndServeFunc: func() error {
			serverWg.Done()
			return nil
		},
	}
	failingServerMock = &mocks_service.HTTPServerMock{
		ListenAndServeFunc: func() error {
			serverWg.Done()
			return errServer
		},
	}
	funcDoGetHTTPServerOK = func(bindAddr string, router http.Handler) service.HTTPServer {
		return serverMock
	}
	funcDoGetHTTPServerFail = func(bindAddr string, router http.Handler) service.HTTPServer {
		return failingServerMock
	}

	// Health Client Mock
	funcDoGetHealthClient = func(name string, url string) *health.Client {
		return &health.Client{
			URL:    url,
			Name:   name,
			Client: newMockHTTPClient(&http.Response{}, nil),
		}
	}

	//Interactive API client
	funcDoGetInteractivesAPIClientOk = func(apiRouter *health.Client) (routes.InteractivesAPIClient, error) {
		return &mocks_routes.InteractivesAPIClientMock{
			ListInteractivesFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, q *interactives.Filter) ([]interactives.Interactive, error) {
				return []interactives.Interactive{}, nil
			},
		}, nil
	}

	funcDoGetStorageProviderOk = func(_ *config.Config, _ storage.DownloadServiceAPIClient) (storage.Provider, error) {
		return &mocks_storage.ProviderMock{
			CheckerFunc: func() func(context.Context, *healthcheck.CheckState) error {
				return func(_ context.Context, s *healthcheck.CheckState) error {
					return s.Update(healthcheck.StatusOK, "mocked storage provider healthy", 0)
				}
			},
			GetFunc: func(context.Context, string) (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader("content")), nil
			},
		}, nil
	}

	funcDoGetDownloadServiceAPIClientOk = func(_ *config.Config) (storage.DownloadServiceAPIClient, error) {
		return &mocks_storage.DownloadServiceAPIClientMock{
			CheckerFunc: func(_ context.Context, s *healthcheck.CheckState) error {
				return s.Update(healthcheck.StatusOK, "mocked storage provider healthy", 0)
			},
			DownloadFunc: func(context.Context, string) (*download.Response, error) {
				return &download.Response{Content: io.NopCloser(strings.NewReader("content"))}, nil
			},
		}, nil
	}
)

// NewMockHTTPClient mocks HTTP Client
func newMockHTTPClient(r *http.Response, err error) *dphttp.ClienterMock {
	return &dphttp.ClienterMock{
		SetPathsWithNoRetriesFunc: func(paths []string) {},
		GetPathsWithNoRetriesFunc: func() []string { return []string{} },
		DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return r, err
		},
	}
}

func TestConstructorNew(t *testing.T) {
	Convey("New returns a new uninitialised service", t, func() {
		So(service.New(nil, nil), ShouldResemble, &service.Service{})
	})
	Convey("Then service is initialised successfully", t, func() {
		initMock := &mocks_service.InitialiserMock{
			DoGetHealthClientFunc:          funcDoGetHealthClient,
			DoGetHealthCheckFunc:           funcDoGetHealthCheckOK,
			DoGetHTTPServerFunc:            funcDoGetHTTPServerOK,
			DoGetInteractivesAPIClientFunc: funcDoGetInteractivesAPIClientOk,
			DoGetStorageProviderFunc:       funcDoGetStorageProviderOk,
		}
		mockServiceList := service.NewServiceList(initMock)
		cfg, err := config.Get()

		svc := service.New(cfg, mockServiceList)

		So(err, ShouldBeNil)
		So(svc.Config, ShouldResemble, cfg)
		So(svc.ServiceList, ShouldResemble, mockServiceList)
	})
}

func TestInitSuccess(t *testing.T) {
	Convey("Given all dependencies are successfully initialised", t, func() {
		initMock := &mocks_service.InitialiserMock{
			DoGetHealthClientFunc:             funcDoGetHealthClient,
			DoGetHealthCheckFunc:              funcDoGetHealthCheckOK,
			DoGetHTTPServerFunc:               funcDoGetHTTPServerOK,
			DoGetInteractivesAPIClientFunc:    funcDoGetInteractivesAPIClientOk,
			DoGetStorageProviderFunc:          funcDoGetStorageProviderOk,
			DoGetDownloadServiceAPIClientFunc: funcDoGetDownloadServiceAPIClientOk,
		}
		mockServiceList := service.NewServiceList(initMock)

		Convey("and valid config and service error channel are provided", func() {
			service.BuildTime = "TestBuildTime"
			service.GitCommit = "TestGitCommit"
			service.Version = "TestVersion"

			cfg, _ := config.Get()
			svc := service.New(cfg, mockServiceList)

			Convey("When Init is called", func() {
				err := svc.Init(ctx)

				Convey("Then service is initialised successfully", func() {
					So(svc.HealthCheck, ShouldResemble, hcMock)
					So(svc.Server, ShouldResemble, serverMock)
					//todo add handlers check?

					Convey("And returns no errors", func() {
						So(err, ShouldBeNil)

						Convey("And the checkers are registered and the healthcheck", func() {
							So(mockServiceList.HealthCheck, ShouldBeTrue)
							So(len(hcMock.AddCheckCalls()), ShouldEqual, expectedCheckCalls)
							So(hcMock.AddCheckCalls()[0].Name, ShouldResemble, "storage provider")
							So(hcMock.AddCheckCalls()[1].Name, ShouldResemble, "API router")
							So(len(initMock.DoGetHTTPServerCalls()), ShouldEqual, 1)
							So(initMock.DoGetHTTPServerCalls()[0].BindAddr, ShouldEqual, ":27300")
						})
					})
				})
			})
		})
	})
}

func TestInitFailure(t *testing.T) {
	Convey("Given failure to create healthcheck", t, func() {
		initMock := &mocks_service.InitialiserMock{
			DoGetHealthClientFunc:             funcDoGetHealthClient,
			DoGetHealthCheckFunc:              funcDoGetHealthCheckFail,
			DoGetInteractivesAPIClientFunc:    funcDoGetInteractivesAPIClientOk,
			DoGetStorageProviderFunc:          funcDoGetStorageProviderOk,
			DoGetDownloadServiceAPIClientFunc: funcDoGetDownloadServiceAPIClientOk,
		}
		mockServiceList := service.NewServiceList(initMock)

		Convey("and valid config and service error channel are provided", func() {
			service.BuildTime = "TestBuildTime"
			service.GitCommit = "TestGitCommit"
			service.Version = "TestVersion"

			cfg, _ := config.Get()
			svc := service.New(cfg, mockServiceList)

			Convey("When Init is called", func() {
				err := svc.Init(ctx)

				Convey("Then service initialisation fails", func() {
					So(svc.ServiceList.HealthCheck, ShouldBeFalse)

					// Healthcheck and Server not initialised
					So(svc.HealthCheck, ShouldBeNil)
					So(svc.Server, ShouldBeNil)

					Convey("And returns error", func() {
						So(err, ShouldNotBeNil)
						So(errors.Unwrap(err), ShouldResemble, errHealthCheck)
					})
				})
			})
		})
	})

	Convey("Given that Checkers cannot be registered", t, func() {
		initMock := &mocks_service.InitialiserMock{
			DoGetHealthClientFunc:             funcDoGetHealthClient,
			DoGetHealthCheckFunc:              funcDoGetHealthAddCheckerFail,
			DoGetInteractivesAPIClientFunc:    funcDoGetInteractivesAPIClientOk,
			DoGetStorageProviderFunc:          funcDoGetStorageProviderOk,
			DoGetDownloadServiceAPIClientFunc: funcDoGetDownloadServiceAPIClientOk,
		}
		mockServiceList := service.NewServiceList(initMock)

		Convey("and valid config and service error channel are provided", func() {
			service.BuildTime = "TestBuildTime"
			service.GitCommit = "TestGitCommit"
			service.Version = "TestVersion"

			cfg, _ := config.Get()
			svc := service.New(cfg, mockServiceList)

			Convey("When Init is called", func() {
				err := svc.Init(ctx)

				Convey("Then service initialisation fails", func() {
					So(svc.ServiceList.HealthCheck, ShouldBeTrue)
					So(svc.HealthCheck, ShouldResemble, hcMockAddFail)

					// Server not initialised
					So(svc.Server, ShouldBeNil)

					Convey("And returns error", func() {
						So(err, ShouldNotBeNil)
						So(errors.Unwrap(err), ShouldResemble, errAddCheckFail)

						Convey("And all checks try to register", func() {
							So(mockServiceList.HealthCheck, ShouldBeTrue)
							So(len(hcMockAddFail.AddCheckCalls()), ShouldEqual, expectedCheckCalls)
						})
					})
				})
			})
		})
	})
}

func TestStart(t *testing.T) {
	Convey("Given a correctly initialised Service with mocked dependencies", t, func() {
		initMock := &mocks_service.InitialiserMock{
			DoGetHealthClientFunc:             funcDoGetHealthClient,
			DoGetHealthCheckFunc:              funcDoGetHealthCheckOK,
			DoGetHTTPServerFunc:               funcDoGetHTTPServerOK,
			DoGetInteractivesAPIClientFunc:    funcDoGetInteractivesAPIClientOk,
			DoGetStorageProviderFunc:          funcDoGetStorageProviderOk,
			DoGetDownloadServiceAPIClientFunc: funcDoGetDownloadServiceAPIClientOk,
		}
		serverWg.Add(1)

		cfg, _ := config.Get()
		mockServiceList := service.NewServiceList(initMock)
		svc := &service.Service{
			Config:      cfg,
			HealthCheck: hcMock,
			Server:      serverMock,
			ServiceList: mockServiceList,
		}

		svcErrors := make(chan error, 1)

		Convey("When service starts", func() {
			svc.Run(ctx, svcErrors)

			Convey("Then healthcheck is started and HTTP server starts listening", func() {
				So(len(hcMock.StartCalls()), ShouldEqual, 1)
				serverWg.Wait() // Wait for HTTP server go-routine to finish
				So(len(serverMock.ListenAndServeCalls()), ShouldEqual, 1)
			})
		})
	})

	Convey("Given that HTTP Server fails", t, func() {
		initMock := &mocks_service.InitialiserMock{
			DoGetHealthClientFunc: funcDoGetHealthClient,
			DoGetHealthCheckFunc:  funcDoGetHealthCheckOK,
			DoGetHTTPServerFunc:   funcDoGetHTTPServerFail,
		}
		serverWg.Add(1)

		Convey("and valid config and service error channel are provided", func() {
			service.BuildTime = "TestBuildTime"
			service.GitCommit = "TestGitCommit"
			service.Version = "TestVersion"

			cfg, _ := config.Get()
			mockServiceList := service.NewServiceList(initMock)
			svc := &service.Service{
				Config:      cfg,
				HealthCheck: hcMock,
				Server:      failingServerMock,
				ServiceList: mockServiceList,
			}

			svcErrors := make(chan error, 1)

			Convey("When service starts", func() {
				svc.Run(ctx, svcErrors)

				Convey("Then service start fails and returns an error in the error channel", func() {
					sErr := <-svcErrors
					So(errors.Unwrap(sErr), ShouldResemble, errServer)
					So(len(failingServerMock.ListenAndServeCalls()), ShouldEqual, 1)
				})
			})
		})
	})
}

func TestCloseSuccess(t *testing.T) {
	Convey("Given a correctly initialised service", t, func() {
		hcStopped := false

		// healthcheck Stop does not depend on any other service being closed/stopped
		hcCloseMock := &mocks_service.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
			StopFunc:     func() { hcStopped = true },
		}

		// server Shutdown will fail if healthcheck is not stopped
		serverCloseMock := &mocks_service.HTTPServerMock{
			ListenAndServeFunc: func() error { return nil },
			ShutdownFunc: func(ctx context.Context) error {
				if !hcStopped {
					return errors.New("Server stopped before healthcheck")
				}
				return nil
			},
		}

		cfg, _ := config.Get()
		serviceList := service.NewServiceList(nil)
		serviceList.HealthCheck = true
		svc := service.Service{
			Config:      cfg,
			HealthCheck: hcCloseMock,
			Server:      serverCloseMock,
			ServiceList: serviceList,
		}

		Convey("When closing service", func() {
			err := svc.Close(ctx)

			Convey("Then it results in all the dependencies being closed in the expected order", func() {
				So(err, ShouldBeNil)
				So(len(hcCloseMock.StopCalls()), ShouldEqual, 1)
				So(len(serverCloseMock.ShutdownCalls()), ShouldEqual, 1)
			})
		})
	})
}

func TestCloseFailure(t *testing.T) {
	Convey("Given if service fails to stop", t, func() {
		failingServerCloseMock := &mocks_service.HTTPServerMock{
			ListenAndServeFunc: func() error { return nil },
			ShutdownFunc: func(ctx context.Context) error {
				return errors.New("Failed to stop http server")
			},
		}

		Convey("And given a correctly initialised service", func() {
			cfg, _ := config.Get()
			serviceList := service.NewServiceList(nil)
			serviceList.HealthCheck = true
			svc := service.Service{
				Config:      cfg,
				HealthCheck: hcMock,
				Server:      failingServerCloseMock,
				ServiceList: serviceList,
			}

			Convey("When closing the service", func() {
				err := svc.Close(ctx)

				Convey("Then Close operation tries to close all dependencies", func() {
					So(len(hcMock.StopCalls()), ShouldEqual, 1)
					So(len(failingServerCloseMock.ShutdownCalls()), ShouldEqual, 1)

					Convey("And returns an error", func() {
						So(err, ShouldNotBeNil)
						So(err.Error(), ShouldResemble, "failed to shutdown gracefully")
					})
				})
			})
		})
	})

	Convey("Given that a dependency takes more time to close than the graceful shutdown timeout", t, func() {
		hcStopped := false

		// healthcheck Stop does not depend on any other service being closed/stopped
		hcShutdownCloseMock := &mocks_service.HealthCheckerMock{
			StopFunc: func() { hcStopped = true },
		}

		// server Shutdown will fail if healthcheck is not stopped
		serverShutdownCloseMock := &mocks_service.HTTPServerMock{
			ShutdownFunc: func(ctx context.Context) error {
				if !hcStopped {
					return errors.New("Server was stopped before healthcheck")
				}
				return nil
			},
		}

		serverShutdownCloseMock.ShutdownFunc = func(ctx context.Context) error {
			time.Sleep(20 * time.Millisecond)
			return nil
		}

		Convey("And given a correctly initialised service", func() {
			cfg, _ := config.Get()
			cfg.GracefulShutdownTimeout = 1 * time.Millisecond
			serviceList := service.NewServiceList(nil)
			serviceList.HealthCheck = true
			svc := service.Service{
				Config:      cfg,
				HealthCheck: hcShutdownCloseMock,
				Server:      serverShutdownCloseMock,
				ServiceList: serviceList,
			}

			Convey("When closing the service", func() {
				err := svc.Close(ctx)

				Convey("Then closing the service fails with context.DeadlineExceeded error and no further dependencies are attempted to close", func() {
					So(err, ShouldResemble, context.DeadlineExceeded)
					So(len(hcShutdownCloseMock.StopCalls()), ShouldEqual, 1)
					So(len(serverShutdownCloseMock.ShutdownCalls()), ShouldEqual, 1)
				})
			})
		})
	})
}

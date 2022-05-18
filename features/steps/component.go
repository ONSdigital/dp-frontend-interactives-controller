package steps

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-api-clients-go/v2/interactives"
	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-frontend-interactives-controller/config"
	"github.com/ONSdigital/dp-frontend-interactives-controller/routes"
	mocks_routes "github.com/ONSdigital/dp-frontend-interactives-controller/routes/mocks"
	"github.com/ONSdigital/dp-frontend-interactives-controller/service"
	mocks_service "github.com/ONSdigital/dp-frontend-interactives-controller/service/mocks"
	"github.com/ONSdigital/dp-frontend-interactives-controller/storage"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dplog "github.com/ONSdigital/log.go/v2/log"
	"github.com/chromedp/chromedp"
	"github.com/maxcnunes/httpfake"
)

type Chrome struct {
	execAllocatorCanceller context.CancelFunc
	ctxCanceller           context.CancelFunc
	ctx                    context.Context
}

type Component struct {
	componenttest.ErrorFeature
	svc         *service.Service
	errorChan   chan error
	HTTPServer  *http.Server
	ctx         context.Context
	chrome      Chrome
	fakeRequest *httpfake.Request
}

func NewComponent() (*Component, error) {
	c := &Component{
		HTTPServer: &http.Server{},
		errorChan:  make(chan error),
		ctx:        context.Background(),
	}

	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	cfg.ServeFromEmbeddedContent = true

	initFunctions := &mocks_service.InitialiserMock{
		DoGetHTTPServerFunc:            c.DoGetHTTPServer,
		DoGetHealthCheckFunc:           DoGetHealthcheckOk,
		DoGetHealthClientFunc:          DoGetHealthClient,
		DoGetInteractivesAPIClientFunc: DoGetInteractivesAPIClientFunc,
		DoGetStorageProviderFunc:       DoGetStorageProvider,
	}

	serviceList := service.NewServiceList(initFunctions)

	c.runApplication(cfg, serviceList)

	return c, nil
}

func (c *Component) Reset() *Component {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		// set this to false to be able to watch the browser in action
		chromedp.Flag("headless", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	c.chrome.execAllocatorCanceller = cancel

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	c.chrome.ctxCanceller = cancel

	log.Print("re-starting chrome ...")

	c.chrome.ctx = ctx

	return c
}

func (c *Component) Close() error {
	dplog.Info(c.ctx, "Shutting down app from test ...")
	if c.svc != nil {
		_ = c.svc.Close(c.ctx)
	}

	c.chrome.ctxCanceller()
	c.chrome.execAllocatorCanceller()

	return nil
}

func (c *Component) InitialiseService() (http.Handler, error) {
	return c.HTTPServer.Handler, nil
}

func (c *Component) DoGetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	c.HTTPServer.Addr = bindAddr
	c.HTTPServer.Handler = router
	return c.HTTPServer
}

func (c *Component) runApplication(cfg *config.Config, svcList *service.ExternalServiceList) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	svc := service.New(cfg, svcList)
	if err := svc.Init(c.ctx); err != nil {
		dplog.Error(c.ctx, "failed to initialise service", err)
		return
	}

	go func() {
		svc.Run(c.ctx, c.errorChan)

		// blocks until an os interrupt or a fatal error occurs
		select {
		case err := <-c.errorChan:
			dplog.Error(c.ctx, "service error received", err)
		case sig := <-signals:
			dplog.Info(c.ctx, "os signal received", dplog.Data{"signal": sig})
		}
	}()
}

// DoGetHealthCheck creates a healthcheck with versionInfo
func DoGetHealthcheckOk(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
	return &mocks_service.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
	}, nil
}

// GetHealthClient returns a healthclient for the provided URL
func DoGetHealthClient(name, url string) *health.Client {
	return &health.Client{}
}

func DoGetInteractivesAPIClientFunc(apiRouter *health.Client) (routes.InteractivesAPIClient, error) {
	return &mocks_routes.InteractivesAPIClientMock{
		ListInteractivesFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, q *interactives.InteractiveFilter) ([]interactives.Interactive, error) {
			pub := true
			return []interactives.Interactive{
				{
					ID:        "123456",
					Published: &pub,
					Metadata: &interactives.InteractiveMetadata{
						HumanReadableSlug: "a-slug", ResourceID: "resid123",
					},
					Archive: &interactives.InteractiveArchive{
						Files: []*interactives.InteractiveFile{
							{Name: "/index.html"},
						},
					},
				},
			}, nil
		},
	}, nil
}

func DoGetStorageProvider(_ *config.Config, _ storage.DownloadServiceAPIClient) (storage.Provider, error) {
	//note: real implementation (mock is available: mocks_storage.ProviderMock{})
	return storage.NewFromEmbeddedFilesystem(), nil
}

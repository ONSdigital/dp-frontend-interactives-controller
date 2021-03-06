package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/dp-frontend-interactives-controller/config"
	"github.com/ONSdigital/dp-frontend-interactives-controller/service"
	"github.com/ONSdigital/log.go/v2/log"
)

func main() {
	log.Namespace = "dp-frontend-interactives-controller"
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(ctx, "application unexpectedly failed", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run(ctx context.Context) error {
	// Create error channel for os signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Read config
	cfg, err := config.Get()
	if err != nil {
		log.Fatal(ctx, "failed to retrieve service configuration", err)
		return err
	}

	log.Info(ctx, "got service configuration", log.Data{"config": cfg})

	// Create service initialiser and an error channel for service errors
	svcList := service.NewServiceList(&service.Init{})
	svcErrors := make(chan error, 1)

	// Run service
	svc := service.New(cfg, svcList)
	if err := svc.Init(ctx); err != nil {
		log.Fatal(ctx, "failed to initialise service", err)
		return err
	}
	svc.Run(ctx, svcErrors)

	// Blocks until an os interrupt or a fatal error occurs
	select {
	case err := <-svcErrors:
		log.Error(ctx, "service error received", err)
	case sig := <-signals:
		log.Info(ctx, "os signal received", log.Data{"signal": sig})
	}

	return svc.Close(ctx)
}

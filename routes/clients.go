package routes

import (
	"context"
	"github.com/ONSdigital/dp-api-clients-go/v2/interactives"
)

//go:generate moq -out mocks/api_clients.go -pkg mocks_routes . InteractivesAPIClient

type InteractivesAPIClient interface {
	ListInteractives(ctx context.Context, userAuthToken, serviceAuthToken string, f *interactives.InteractiveFilter) ([]interactives.Interactive, error)
}

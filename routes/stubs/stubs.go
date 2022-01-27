package stubs

import "context"

//todo this should come from: https://github.com/ONSdigital/dp-api-clients-go

type InteractivesAPIGetResponse struct {
	slug string
}

// InteractivesApiClient is an interface for the Interactives API client
type InteractivesAPIClient interface {
	Get(ctx context.Context, userAccessToken, collectionID, lang, uri string) (*InteractivesAPIGetResponse, error)
}

type StubbedInteractivesAPIClient struct{}

func (c StubbedInteractivesAPIClient) Get(ctx context.Context, userAccessToken, collectionID, lang, uri string) (*InteractivesAPIGetResponse, error) {
	return &InteractivesAPIGetResponse{}, nil
}

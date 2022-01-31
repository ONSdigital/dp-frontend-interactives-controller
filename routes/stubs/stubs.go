package stubs

import "context"

//todo this should come from: https://github.com/ONSdigital/dp-api-clients-go

type InteractivesAPIGetResponse struct {
	Metadata map[string]string
}

// InteractivesApiClient is an interface for the Interactives API client
type InteractivesAPIClient interface {
	Get(ctx context.Context, id string) (*InteractivesAPIGetResponse, error)
}

type StubbedInteractivesAPIClient struct{}

func (c StubbedInteractivesAPIClient) Get(_ context.Context, id string) (*InteractivesAPIGetResponse, error) {
	return &InteractivesAPIGetResponse{
		Metadata: map[string]string{
			"id":                  id,
			"human_readable_slug": "an_url_slug",
			"path":                "interactives/an_interactives_id",
		},
	}, nil
}

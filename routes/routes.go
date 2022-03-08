package routes

import (
	"fmt"
	"github.com/ONSdigital/dp-frontend-interactives-controller/config"
	"net/http"

	"github.com/ONSdigital/dp-frontend-interactives-controller/routes/stubs"
	"github.com/ONSdigital/dp-frontend-interactives-controller/storage"
	"github.com/gorilla/mux"
)

const (
	HealthEndpoint   = "/health"
	EmbeddedSuffix   = "/embed"
	ResourceIdVarKey = "resource_id"
	SlugVarKey       = "human_readable_slug"
)

// Clients - struct containing all the clients for the controller
type Clients struct {
	Storage storage.Provider
	Api     stubs.InteractivesAPIClient
}

// Setup registers routes for the service
func Setup(cfg *config.Config, r *mux.Router, hc http.HandlerFunc, interactivesHandler http.HandlerFunc) {
	r.StrictSlash(true).Path(HealthEndpoint).HandlerFunc(hc)

	// slug and resource_id (+ /embed)
	r.StrictSlash(true).Path(getPath(false, true)).Methods(http.MethodGet).HandlerFunc(interactivesHandler)
	r.StrictSlash(true).Path(getPath(true, true)).Methods(http.MethodGet).HandlerFunc(interactivesHandler)
	// just resource_id (+ /embed)
	r.StrictSlash(true).Path(getPath(false, false)).Methods(http.MethodGet).HandlerFunc(interactivesHandler)
	r.StrictSlash(true).Path(getPath(true, false)).Methods(http.MethodGet).HandlerFunc(interactivesHandler)

	if cfg.ServeFromEmbeddedContent {
		r.StrictSlash(true).PathPrefix("/interactives").Methods(http.MethodGet).HandlerFunc(interactivesHandler)
	}
}

func getPath(withEmbed, withSlug bool) string {
	resourceTypeKey := "interactives" //this is driven from dp-frontend-router (should be 'interactives')

	resourceIdPattern := "[a-zA-Z0-9]{8}"
	url := fmt.Sprintf("/{%s}/{%s:%s}", resourceTypeKey, ResourceIdVarKey, resourceIdPattern)
	if withSlug {
		slugKeyPattern := "[a-zA-Z0-9\\-]+"
		url = fmt.Sprintf("/{%s}/{%s:%s}-{%s:%s}", resourceTypeKey, SlugVarKey, slugKeyPattern, ResourceIdVarKey, resourceIdPattern)
	}
	if withEmbed {
		url = fmt.Sprintf("%s%s", url, EmbeddedSuffix)
	}

	return url
}

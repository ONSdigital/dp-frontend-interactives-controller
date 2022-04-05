package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-frontend-interactives-controller/config"
	"github.com/ONSdigital/dp-frontend-interactives-controller/storage"
	"github.com/gorilla/mux"
)

const (
	ResourceTypeKey  = "interactives"
	HealthEndpoint   = "/health"
	EmbeddedSuffix   = "/embed"
	ResourceIdVarKey = "resource_id"
	SlugVarKey       = "human_readable_slug"
	CatchAllVarKey   = "uri"
)

// Clients - struct containing all the clients for the controller
type Clients struct {
	Storage storage.Provider
	API     InteractivesAPIClient
}

// Setup registers routes for the service
func Setup(_ *config.Config, r *mux.Router, hc http.HandlerFunc, interactivesHandler http.HandlerFunc, redirectHandler http.HandlerFunc) {
	r.StrictSlash(true).Path(HealthEndpoint).HandlerFunc(hc)
	// /interactives
	r.StrictSlash(true).
		PathPrefix(getPath(false, true)).
		MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
			return !strings.HasSuffix(r.URL.Path, EmbeddedSuffix)
		}).
		Methods(http.MethodGet).
		Handler(interactivesHandler)

	// only resource_id - redirect
	r.StrictSlash(true).
		PathPrefix(getPath(false, false)).
		MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
			return !strings.HasSuffix(r.URL.Path, EmbeddedSuffix)
		}).
		Methods(http.MethodGet).
		Handler(redirectHandler)
	r.StrictSlash(true).
		Path(getPath(true, false)).
		Methods(http.MethodGet).
		HandlerFunc(redirectHandler)

	// fixed /embed URLs
	r.StrictSlash(true).
		Path(getPath(true, true)).
		Methods(http.MethodGet).
		HandlerFunc(interactivesHandler)
}

func getPath(withEmbed, withSlug bool) string {

	resourceIdPattern := "[a-zA-Z0-9]{8}"
	url := fmt.Sprintf("/{%s}/{%s:%s}", ResourceTypeKey, ResourceIdVarKey, resourceIdPattern)
	if withSlug {
		slugKeyPattern := "[a-zA-Z0-9\\-]+"
		url = fmt.Sprintf("/{%s}/{%s:%s}-{%s:%s}", ResourceTypeKey, SlugVarKey, slugKeyPattern, ResourceIdVarKey, resourceIdPattern)
	}
	if withEmbed {
		url = fmt.Sprintf("%s%s", url, EmbeddedSuffix)
	} else {
		url = fmt.Sprintf("%s{%s:.*}", url, CatchAllVarKey)
	}

	return url
}

package routes

import (
	"context"
	"net/http"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

const (
	HealthEndpoint       = "/health"
	InteractivesEndpoint = "/interactives"
)

// Clients - struct containing all the clients for the controller
type Clients struct {
	HealthCheckHandler  func(w http.ResponseWriter, req *http.Request)
	InteractivesHandler func(w http.ResponseWriter, req *http.Request)
}

// Setup registers routes for the service
func Setup(ctx context.Context, r *mux.Router, c Clients) {
	log.Info(ctx, "adding routes")
	r.StrictSlash(true).Path(HealthEndpoint).HandlerFunc(c.HealthCheckHandler)
	r.StrictSlash(true).PathPrefix(InteractivesEndpoint).HandlerFunc(c.InteractivesHandler)
}

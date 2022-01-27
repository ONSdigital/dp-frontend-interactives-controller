package routes

import (
	"github.com/ONSdigital/dp-frontend-interactives-controller/routes/stubs"
	"github.com/ONSdigital/dp-frontend-interactives-controller/storage"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	HealthEndpoint = "/health"
)

// Clients - struct containing all the clients for the controller
type Clients struct {
	Storage storage.Provider
	Api     stubs.InteractivesAPIClient
}

// Setup registers routes for the service
func Setup(r *mux.Router, hc http.HandlerFunc, interactivesHandler http.HandlerFunc) {
	r.StrictSlash(true).Path(HealthEndpoint).HandlerFunc(hc)
	r.StrictSlash(true).Path("/{uri:.*}").Methods("GET").HandlerFunc(interactivesHandler)
}

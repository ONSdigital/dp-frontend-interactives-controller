package routes_test

import (
	"context"
	"github.com/ONSdigital/dp-frontend-interactives-controller/routes"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	statusNoContentFunc = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
)

func TestRoutes(t *testing.T) {

	Convey("Given clients setup to return StatusNoContent[204] for all handlers", t, func() {
		clients := routes.Clients{HealthCheckHandler: statusNoContentFunc, InteractivesHandler: statusNoContentFunc}
		r := mux.NewRouter()
		routes.Setup(context.TODO(), r, clients)

		Convey("when "+routes.HealthEndpoint+" is called", func() {
			req := httptest.NewRequest("GET", routes.HealthEndpoint, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			Convey("then 204 is returned", func() {
				So(w.Code, ShouldEqual, http.StatusNoContent)
			})
		})

		Convey("when "+routes.InteractivesEndpoint+" is called", func() {
			req := httptest.NewRequest("GET", routes.InteractivesEndpoint, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			Convey("then 204 is returned", func() {
				So(w.Code, ShouldEqual, http.StatusNoContent)
			})
		})

		Convey("when an unknown endpoint is called", func() {
			req := httptest.NewRequest("GET", "/endpoint-does-not-exist", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			Convey("then 404 is returned", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}

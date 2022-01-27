package routes_test

import (
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
	Convey("Given router setup to return StatusNoContent[204] for healthcheck and interactives", t, func() {

		r := mux.NewRouter()
		routes.Setup(r, statusNoContentFunc, statusNoContentFunc)

		Convey("when "+routes.HealthEndpoint+" is called", func() {
			req := httptest.NewRequest("GET", routes.HealthEndpoint, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			Convey("then 204 is returned", func() {
				So(w.Code, ShouldEqual, http.StatusNoContent)
			})
		})

		Convey("when a GET is called on any path", func() {
			req := httptest.NewRequest("GET", "/method-supported", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			Convey("then 204 is returned", func() {
				So(w.Code, ShouldEqual, http.StatusNoContent)
			})
		})

		Convey("when another method is called", func() {
			req := httptest.NewRequest("POST", "/method-not-supported", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			Convey("then 405 is returned", func() {
				So(w.Code, ShouldEqual, http.StatusMethodNotAllowed)
			})
		})
	})
}

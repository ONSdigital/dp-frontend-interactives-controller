package routes_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-frontend-interactives-controller/routes"
	"github.com/gorilla/mux"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	statusNoContentFunc = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
)

func checkPathVariablesHandler(t *testing.T, slug, resourceId string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		actualResourceId := vars[routes.ResourceIdVarKey]
		actualSlug := vars[routes.SlugVarKey]

		if actualResourceId != resourceId {
			t.Errorf("resourceId not as expected, expected=[%s], actual=[%s]", resourceId, actualResourceId)
		}
		if actualSlug != slug {
			t.Errorf("slug not as expected, expected=[%s], actual=[%s]", slug, actualSlug)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func TestSetup(t *testing.T) {
	Convey("Given setup then 5 routes are applied", t, func() {
		r := mux.NewRouter()
		routes.Setup(r, statusNoContentFunc, statusNoContentFunc)

		routes := 0
		err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			_, err := route.GetPathTemplate()
			if err != nil {
				return err
			}
			routes++
			return nil
		})

		So(err, ShouldBeNil)
		So(routes, ShouldEqual, 6) //6 with static
	})
}

func TestRoutes(t *testing.T) {
	resourceType := "interactives"
	validResourceId := "abcde123"
	validSlug := "nice-readable-slug"

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

		Convey("when a mapped route is called with trailing slash", func() {
			urlWithTrailingSlash := fmt.Sprintf("/%s/%s-%s/", resourceType, validSlug, validResourceId)
			req := httptest.NewRequest(http.MethodGet, urlWithTrailingSlash, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			Convey("then 301 is returned for nonslash url", func() {
				So(w.Header().Get("location"), ShouldEqual, urlWithTrailingSlash[:len(urlWithTrailingSlash)-1])
				So(w.Code, ShouldEqual, http.StatusMovedPermanently)
			})
		})
	})

	type test struct{ method, url, slug, resourceId string }

	Convey("Given router setup to check route variables", t, func() {
		Convey("when a mapped route is called", func() {
			cases := map[string]test{
				"slug-and-resource-id":          {http.MethodGet, fmt.Sprintf("/%s/%s-%s", resourceType, validSlug, validResourceId), validSlug, validResourceId},
				"embedded-slug-and-resource-id": {http.MethodGet, fmt.Sprintf("/%s/%s-%s/embed", resourceType, validSlug, validResourceId), validSlug, validResourceId},
				"resource-id":                   {http.MethodGet, fmt.Sprintf("/%s/%s", resourceType, validResourceId), "", validResourceId},
				"embedded-resource-id":          {http.MethodGet, fmt.Sprintf("/%s/%s/embed", resourceType, validResourceId), "", validResourceId},
			}

			for name, testReq := range cases {

				req := httptest.NewRequest(testReq.method, testReq.url, nil)
				w := httptest.NewRecorder()
		
				h := checkPathVariablesHandler(t, testReq.slug, testReq.resourceId)
				r := mux.NewRouter()
				routes.Setup(r, nil, h)

				r.ServeHTTP(w, req)

				Convey(fmt.Sprintf("then 204 is returned for %s", name), func() {
					So(w.Code, ShouldEqual, http.StatusNoContent)
				})
			}
		})

		Convey("when am unsupported route is called", func() {
			cases := map[string]test{
				"not-supported":     {http.MethodGet, "/not-supported", "", ""},
				"bad_slug_1":        {http.MethodGet, fmt.Sprintf("/%s/%s-%s", resourceType, "", validResourceId), "", ""},
				"bad_slug_2":        {http.MethodGet, fmt.Sprintf("/%s/%s-%s", resourceType, "under_score", validResourceId), "", ""},
				"bad_slug_3":        {http.MethodGet, fmt.Sprintf("/%s/%s-%s", resourceType, "full.stop", validResourceId), "", ""},
				"bad_resource_id_1": {http.MethodGet, fmt.Sprintf("/%s/%s-%s", resourceType, validSlug, ""), "", ""},
				"bad_resource_id_2": {http.MethodGet, fmt.Sprintf("/%s/%s-%s", resourceType, validSlug, "abcde"), "", ""},
				"bad_resource_id_3": {http.MethodGet, fmt.Sprintf("/%s/%s-%s", resourceType, validSlug, "abc-de"), "", ""},
			}

			for name, testReq := range cases {
				req := httptest.NewRequest(testReq.method, testReq.url, nil)
				w := httptest.NewRecorder()

				h := checkPathVariablesHandler(t, testReq.slug, testReq.resourceId)
				r := mux.NewRouter()
				routes.Setup(r, nil, h)
				r.ServeHTTP(w, req)

				Convey(fmt.Sprintf("then 404 is returned for %s", name), func() {
					So(w.Code, ShouldEqual, http.StatusNotFound)
				})
			}
		})
	})
}

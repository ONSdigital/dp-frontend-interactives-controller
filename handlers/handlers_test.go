package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/interactives"
	"github.com/ONSdigital/dp-frontend-interactives-controller/config"
	mocks_routes "github.com/ONSdigital/dp-frontend-interactives-controller/routes/mocks"
	"github.com/gorilla/mux"

	"github.com/ONSdigital/dp-frontend-interactives-controller/routes"
	"github.com/ONSdigital/dp-frontend-interactives-controller/storage"
	mocks_storage "github.com/ONSdigital/dp-frontend-interactives-controller/storage/mocks"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	fileOkContent = "some arbitrary bucket content"
)

type testCliError struct{}

func (e *testCliError) Error() string { return "client error" }
func (e *testCliError) Code() int     { return http.StatusNotFound }

func TestSetStatusCode(t *testing.T) {
	Convey("test setStatusCode", t, func() {

		Convey("test status code handles 404 response from client", func() {
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			err := &testCliError{}

			setStatusCode(req, w, http.StatusInternalServerError, err)

			So(w.Code, ShouldEqual, http.StatusNotFound)
		})

		Convey("test status code handles StatusInternalServerError", func() {
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			err := errors.New("internal server error")

			setStatusCode(req, w, http.StatusInternalServerError, err)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})

		Convey("test status code handles StatusNotFound", func() {
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			err := errors.New("not found")

			setStatusCode(req, w, http.StatusNotFound, err)

			So(w.Code, ShouldEqual, http.StatusNotFound)
		})

		Convey("test if status code called twice then not overwritten", func() {
			//https://cs.opensource.google/go/go/+/master:src/net/http/server.go;drc=81431c7aa7c5d782e72dec342442ea7664ef1783;l=141
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			err := errors.New("internal server error")

			setStatusCode(req, w, http.StatusInternalServerError, err)
			setStatusCode(req, w, http.StatusNotFound, err)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
		})
	})
}

func TestInteractives(t *testing.T) {

	Convey("Given a valid s3 bucket and provider", t, func() {
		s3BucketMock := &mocks_storage.S3BucketMock{
			CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error {
				return state.Update(healthcheck.StatusOK, "mocked s3 bucket healthy", 0)
			},
			GetFunc: func(key string) (io.ReadCloser, *int64, error) {
				r := strings.NewReader(fileOkContent)
				contentLen := int64(123)
				return io.NopCloser(r), &contentLen, nil
			},
		}
		storageProvider := storage.NewFromS3Bucket(s3BucketMock)

		Convey("a request to a valid s3 path is made", func() {

			validPath := "/valid/path/to/file.html"
			pub := true
			type test struct {
				expectedStatus, totalCount int
				expectedFileContent, path  string
			}
			cases := map[string]test{
				"happy-path":          {http.StatusOK, 1, fileOkContent, validPath},
				"zero-results":        {http.StatusNotFound, 0, "", validPath},
				"more-than-1-results": {http.StatusNotFound, 2, "", validPath},
			}

			for name, testReq := range cases {
				apiMock := &mocks_routes.InteractivesAPIClientMock{
					ListInteractivesFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, q *interactives.QueryParams) (interactives.List, error) {
						return interactives.List{
							Items: []interactives.Interactive{
								getTestInteractive(pub, nil),
							},
							Count:      1,
							Offset:     0,
							Limit:      10,
							TotalCount: testReq.totalCount,
						}, nil
					},
				}

				clients := routes.Clients{
					Storage: storageProvider,
					API:     apiMock,
				}

				handler := Interactives(&config.Config{}, clients)

				req := httptest.NewRequest(http.MethodGet, testReq.path, nil)
				w := httptest.NewRecorder()
				handler(w, req)

				Convey(fmt.Sprintf("then the status code and body are as expected %s", name), func() {
					res := w.Result()
					defer res.Body.Close()
					data, err := ioutil.ReadAll(res.Body)

					So(err, ShouldBeNil)
					So(w.Code, ShouldEqual, testReq.expectedStatus)
					So(string(data), ShouldEqual, testReq.expectedFileContent)
				})
			}

		})

		Convey("url with only resource-id must redirect", func() {
			pub := true
			mData := &interactives.InteractiveMetadata{HumanReadableSlug: "a-slug", ResourceID: "resid123"}
			apiMock := &mocks_routes.InteractivesAPIClientMock{
				ListInteractivesFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, q *interactives.QueryParams) (interactives.List, error) {
					return interactives.List{
						Items: []interactives.Interactive{
							getTestInteractive(pub, mData),
						},
						Count:      1,
						Offset:     0,
						Limit:      10,
						TotalCount: 1,
					}, nil
				},
			}

			clients := routes.Clients{
				Storage: storageProvider,
				API:     apiMock,
			}

			handler := InteractivesRedirect(&config.Config{}, clients)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			handler(w, req)

			Convey("then the status code is 301", func() {
				expectedRedirect := fmt.Sprintf("/%s-%s%s", mData.HumanReadableSlug, mData.ResourceID, routes.EmbeddedSuffix)
				res := w.Result()
				defer res.Body.Close()
				body := w.Body.String()

				So(res.StatusCode, ShouldEqual, http.StatusMovedPermanently)
				So(strings.Contains(body, expectedRedirect), ShouldBeTrue)
			})
		})

		Convey("mismatched slug must redirect", func() {
			pub := true
			mData := &interactives.InteractiveMetadata{HumanReadableSlug: "a-slug", ResourceID: "resid123"}
			apiMock := &mocks_routes.InteractivesAPIClientMock{
				ListInteractivesFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, q *interactives.QueryParams) (interactives.List, error) {
					return interactives.List{
						Items: []interactives.Interactive{
							getTestInteractive(pub, mData),
						},
						Count:      1,
						Offset:     0,
						Limit:      10,
						TotalCount: 1,
					}, nil
				},
			}

			clients := routes.Clients{
				Storage: storageProvider,
				API:     apiMock,
			}

			handler := Interactives(&config.Config{}, clients)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req = mux.SetURLVars(req, map[string]string{routes.SlugVarKey: "different-slug", routes.ResourceIdVarKey: "resid123"})
			w := httptest.NewRecorder()
			handler(w, req)

			Convey("then the status code is 301", func() {
				expectedRedirect := fmt.Sprintf("/%s-%s%s", mData.HumanReadableSlug, mData.ResourceID, routes.EmbeddedSuffix)
				res := w.Result()
				defer res.Body.Close()
				body := w.Body.String()

				So(res.StatusCode, ShouldEqual, http.StatusMovedPermanently)
				So(strings.Contains(body, expectedRedirect), ShouldBeTrue)
			})
		})
	})
}

func getTestInteractive(published bool, m *interactives.InteractiveMetadata) interactives.Interactive {
	if m == nil {
		m = &interactives.InteractiveMetadata{
			HumanReadableSlug: "slug",
			ResourceID:        "abcd123e",
		}
	}
	return interactives.Interactive{
		ID:        "123456",
		Published: &published,
		Metadata:  m,
		Archive: &interactives.InteractiveArchive{
			Files: []*interactives.InteractiveFile{
				{Name: "/index.html"},
			},
		},
	}
}

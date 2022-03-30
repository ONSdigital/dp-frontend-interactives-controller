package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/ONSdigital/dp-api-clients-go/v2/interactives"
	mocks_routes "github.com/ONSdigital/dp-frontend-interactives-controller/routes/mocks"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
								{ID: "123456", Metadata: nil, Archive: nil},
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

				handler := Interactives(clients)

				req := httptest.NewRequest(http.MethodGet, testReq.path, nil)
				w := httptest.NewRecorder()
				handler(w, req)

				Convey(fmt.Sprintf("then the status code is 200 and body is as expected %s", name), func() {
					res := w.Result()
					defer res.Body.Close()
					data, err := ioutil.ReadAll(res.Body)

					So(err, ShouldBeNil)
					So(w.Code, ShouldEqual, testReq.expectedStatus)
					So(string(data), ShouldEqual, testReq.expectedFileContent)
				})
			}

		})
	})
}

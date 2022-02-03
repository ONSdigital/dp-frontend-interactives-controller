package handlers

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-frontend-interactives-controller/routes"
	"github.com/ONSdigital/dp-frontend-interactives-controller/routes/stubs"
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

			setStatusCode(req, w, err)

			So(w.Code, ShouldEqual, http.StatusNotFound)
		})

		Convey("test status code handles internal server error", func() {
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			err := errors.New("internal server error")

			setStatusCode(req, w, err)

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

			clients := routes.Clients{
				Storage: storageProvider,
				Api:     &stubs.StubbedInteractivesAPIClient{},
			}

			handler := Interactives(clients)

			req := httptest.NewRequest(http.MethodGet, "/valid/path/to/file.html", nil)
			w := httptest.NewRecorder()
			handler(w, req)

			Convey("then the status code is 200 and body is as expected", func() {
				res := w.Result()
				defer res.Body.Close()
				data, err := ioutil.ReadAll(res.Body)

				So(err, ShouldBeNil)
				So(w.Code, ShouldEqual, http.StatusOK)
				So(string(data), ShouldEqual, fileOkContent)
			})
		})
	})
}

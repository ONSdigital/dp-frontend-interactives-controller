package handlers

import (
	"context"
	mocks_storage "github.com/ONSdigital/dp-frontend-interactives-controller/storage/mocks"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	fileOkContent = "some arbitrary bucket content"
)

var (
	s3BucketMock = &mocks_storage.S3BucketMock{
		CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error {
			return state.Update(healthcheck.StatusOK, "mocked s3 bucket healthy", 0)
		},
		GetFunc: func(key string) (io.ReadCloser, *int64, error) {
			r := strings.NewReader(fileOkContent)
			contentLen := int64(123)
			return io.NopCloser(r), &contentLen, nil
		},
	}
)

func TestS3Handlers(t *testing.T) {
	Convey("Given a request to a valid s3 path", t, func() {
		req := httptest.NewRequest(http.MethodGet, "/valid/path/to/file.html", nil)
		w := httptest.NewRecorder()

		s3Handlers := NewS3Backed(s3BucketMock)
		h := s3Handlers.GetInteractivesHandler()

		Convey("a request to a valid s3 path is made", func() {
			h(w, req)

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

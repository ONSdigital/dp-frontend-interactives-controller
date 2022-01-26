package handlers

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/v2/log"
	"net/http"
)

// ClientError is an interface that can be used to retrieve the status code if a client has errored
type ClientError interface {
	Error() string
	Code() int
}

func setStatusCode(r *http.Request, w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if err, ok := err.(ClientError); ok {
		if err.Code() == http.StatusNotFound {
			status = err.Code()
		}
	}
	log.Error(r.Context(), "setting-response-status", err)
	w.WriteHeader(status)
}

type Handlers interface {
	GetInteractivesHandler() func(http.ResponseWriter, *http.Request)
	Checker() func(context.Context, *healthcheck.CheckState) error
}

func NewLocalFilesystemBacked(root http.Dir) localfs {
	return localfs{root: root}
}

type localfs struct {
	root http.Dir
}

func (s localfs) GetInteractivesHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("%s/%s", s.root, r.URL.Path[1:]))
	}
}

func (s localfs) Checker() func(context.Context, *healthcheck.CheckState) error {
	return func(_ context.Context, s *healthcheck.CheckState) error {
		return s.Update(healthcheck.StatusOK, "localfs healthy", 0)
	}
}

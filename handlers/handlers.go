package handlers

import (
	"context"
	"github.com/ONSdigital/dp-frontend-interactives-controller/routes"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"io"
	"mime"
	"net/http"
	"path/filepath"
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

func Interactives(clients routes.Clients) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		streamFromStorageProvider(w, r, clients)
	}
}

func streamFromStorageProvider(w http.ResponseWriter, r *http.Request, clients routes.Clients) {
	ctx := r.Context()

	//TODO get metadata from API
	vars := mux.Vars(r)
	id := vars[routes.ResourceIdVarKey]
	_, err := clients.Api.Get(r.Context(), id)
	if err != nil {
		//todo 404 from error pass back upstream?
		log.Error(ctx, "failed to get from interactives api", err)
		setStatusCode(r, w, err)
		return
	}

	//todo - from interactive
	path := r.URL.Path

	//stream content to response
	var readCloser io.ReadCloser
	readCloser, err = clients.Storage.Get(path)
	if err != nil {
		//todo 404 from error pass back upstream?
		log.Error(ctx, "failed to get stream from storage provider", err)
		setStatusCode(r, w, err)
		return
	}
	defer closeAndLogError(ctx, readCloser)

	//note: has to be before writing body. ref: https://pkg.go.dev/net/http#ResponseWriter.Write
	ctype := mime.TypeByExtension(filepath.Ext(path))
	if ctype != "" {
		w.Header().Set("Content-Type", ctype)
	}

	_, err = io.Copy(w, readCloser)
	if err != nil {
		log.Error(ctx, "failed to write response", err)
		setStatusCode(r, w, err)
		return
	}
}

func closeAndLogError(ctx context.Context, closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Error(ctx, "error closing io.Closer", err)
	}
}

package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
	"path/filepath"

	"github.com/ONSdigital/dp-frontend-interactives-controller/config"

	"github.com/ONSdigital/dp-api-clients-go/v2/interactives"
	"github.com/ONSdigital/dp-frontend-interactives-controller/routes"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// ClientError is an interface that can be used to retrieve the status code if a client has errored
type ClientError interface {
	Error() string
	Code() int
}

func setStatusCode(r *http.Request, w http.ResponseWriter, status int, err error) {
	if e, ok := err.(ClientError); ok {
		if e.Code() == http.StatusNotFound {
			status = e.Code()
		}
	}
	log.Error(r.Context(), "setting-response-status", err)
	w.WriteHeader(status)
}

func Interactives(cfg *config.Config, clients routes.Clients) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		streamFromStorageProvider(w, r, clients, cfg.ServiceAuthToken)
	}
}

func streamFromStorageProvider(w http.ResponseWriter, r *http.Request, clients routes.Clients, serviceAuthToken string) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id := vars[routes.ResourceIdVarKey]

	all, err := clients.API.ListInteractives(r.Context(), "", serviceAuthToken,
		&interactives.QueryParams{
			Offset: 0,
			Limit:  1,
			Filter: &interactives.InteractiveMetadata{ResourceID: id},
		},
	)
	if err != nil {
		setStatusCode(r, w, http.StatusInternalServerError, fmt.Errorf("failed to get from interactives api %w", err))
		return
	}

	if all.TotalCount != 1 {
		setStatusCode(r, w, http.StatusNotFound, fmt.Errorf("cannot find interactive %w", err))
		return
	}

	// block access if interactive is unpublished
	if !*(all.Items[0].Published) {
		setStatusCode(r, w, http.StatusNotFound, errors.New("access prohibited for unpublished interactives"))
	}

	filename := path.Base(r.URL.Path)
	if filename == id || filename == routes.EmbeddedSuffix[1:] { //root url
		filename = "/index.html"
	} else {
		filename = vars[routes.CatchAllVarKey]
	}
	if filename == "" || filename == "/" {
		filename = "/index.html"
	}

	//todo if file not within interactives archive upload files - 404 otherwise get from download svc

	//stream content to response
	var readCloser io.ReadCloser
	readCloser, err = clients.Storage.Get(filename)
	if err != nil {
		//todo 404 from error pass back upstream?
		setStatusCode(r, w, http.StatusInternalServerError, fmt.Errorf("failed to get stream from storage provider %w", err))
		return
	}
	defer closeAndLogError(ctx, readCloser)

	//note: has to be before writing body. ref: https://pkg.go.dev/net/http#ResponseWriter.Write
	ctype := mime.TypeByExtension(filepath.Ext(filename))
	if ctype != "" {
		w.Header().Set("Content-Type", ctype)
	}

	_, err = io.Copy(w, readCloser)
	if err != nil {
		setStatusCode(r, w, http.StatusInternalServerError, fmt.Errorf("failed to write response %w", err))
		return
	}
}

func closeAndLogError(ctx context.Context, closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Error(ctx, "error closing io.Closer", err)
	}
}

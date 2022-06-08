package handlers

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-api-clients-go/v2/interactives"
	"github.com/ONSdigital/dp-frontend-interactives-controller/config"
	"github.com/ONSdigital/dp-frontend-interactives-controller/routes"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"io"
	"mime"
	"net/http"
	"path"
	"path/filepath"
)

const (
	RootFile = "index.html"
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

func InteractivesRedirect(cfg *config.Config, clients routes.Clients) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectToFullyQualifiedURL(w, r, clients, cfg.ServiceAuthToken)
	}
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
	slug := vars[routes.SlugVarKey]

	ix, url := getInteractive(w, r, id, clients, serviceAuthToken)
	if ix == nil {
		setStatusCode(r, w, http.StatusNotFound, fmt.Errorf("failed to find resource id %s", id))
		return
	}

	// redirect if slug mismatch
	if slug != "" && ix.Metadata != nil && ix.Metadata.HumanReadableSlug != slug {
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	}

	filename := path.Base(r.URL.Path)
	if filename == id || filename == routes.EmbeddedSuffix[1:] { //root url
		filename = "/"
	} else {
		filename = vars[routes.CatchAllVarKey][1:] //strip leading /
	}

	var err error
	filename, err = findFile(filename, ix)
	if err != nil {
		setStatusCode(r, w, http.StatusNotFound, fmt.Errorf("cannot find interactive %w", err))
		return
	}

	//stream content to response
	readCloser, err := clients.Storage.Get(r.Context(), filename)
	if err != nil {
		//todo 404 from error pass back upstream? this could be auth - so 404
		setStatusCode(r, w, http.StatusNotFound, fmt.Errorf("failed to get stream from storage provider %w", err))
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

func findFile(filename string, ix *interactives.Interactive) (string, error) {
	if filename == "" || filename == "/" {
		filename = RootFile
	}

	if ix.Archive != nil {
		for _, f := range ix.Archive.Files {
			if f.URI == filename {
				return f.Name, nil
			}
		}
	}

	return "", fmt.Errorf("cannot find root index.html file for %s", ix.ID)
}

func getInteractive(w http.ResponseWriter, r *http.Request, id string, clients routes.Clients, serviceAuthToken string) (*interactives.Interactive, string) {
	all, err := clients.API.ListInteractives(r.Context(), "", serviceAuthToken,
		&interactives.InteractiveFilter{Metadata: &interactives.InteractiveMetadata{ResourceID: id}},
	)
	if err != nil {
		setStatusCode(r, w, http.StatusInternalServerError, fmt.Errorf("failed to get from interactives api %w", err))
		return nil, ""
	}

	if len(all) != 1 {
		setStatusCode(r, w, http.StatusNotFound, fmt.Errorf("cannot find interactive %w", err))
		return nil, ""
	}

	first := &all[0]
	return first, fmt.Sprintf("/%s/%s-%s%s", routes.ResourceTypeKey, first.Metadata.HumanReadableSlug, first.Metadata.ResourceID, routes.EmbeddedSuffix)

}

func redirectToFullyQualifiedURL(w http.ResponseWriter, r *http.Request, clients routes.Clients, serviceAuthToken string) {
	vars := mux.Vars(r)
	id := vars[routes.ResourceIdVarKey]
	ix, url := getInteractive(w, r, id, clients, serviceAuthToken)
	if ix == nil {
		setStatusCode(r, w, http.StatusNotFound, fmt.Errorf("failed to find resource id %s", id))
		return
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func closeAndLogError(ctx context.Context, closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Error(ctx, "error closing io.Closer", err)
	}
}

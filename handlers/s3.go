package handlers

import (
	"context"
	"github.com/ONSdigital/dp-frontend-interactives-controller/storage"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/v2/log"
	"io"
	"mime"
	"net/http"
	"path/filepath"
)

func NewS3Backed(s3bucket storage.S3Bucket) s3 {
	return s3{s3bucket: s3bucket}
}

type s3 struct {
	s3bucket storage.S3Bucket
}

func (s s3) GetInteractivesHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		streamFromS3(w, r, s.s3bucket)
	}
}

func (s s3) Checker() func(context.Context, *healthcheck.CheckState) error {
	return s.s3bucket.Checker
}

func streamFromS3(w http.ResponseWriter, r *http.Request, s3bucket storage.S3Bucket) {
	ctx := r.Context()

	//TODO s3Path from api [for testing bucket/id/./...]
	s3Path := r.URL.Path

	//stream s3 content to response
	var s3ReadCloser io.ReadCloser
	s3ReadCloser, _, err := s3bucket.Get(s3Path)
	if err != nil {
		//todo 404 from error pass back upstream?
		log.Error(ctx, "failed to get stream object from S3 client", err)
		setStatusCode(r, w, err)
		return
	}
	defer closeAndLogError(ctx, s3ReadCloser)

	//note: has to be before writing body. ref: https://pkg.go.dev/net/http#ResponseWriter.Write
	ctype := mime.TypeByExtension(filepath.Ext(s3Path))
	if ctype != "" {
		w.Header().Set("Content-Type", ctype)
	}

	_, err = io.Copy(w, s3ReadCloser)
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

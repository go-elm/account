package service

import (
	"encoding/json"
	"errors"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"

	"golang.org/x/net/context"
)

var (
	// errBadRoute is used for mux errors
	errBadRoute = errors.New("bad route")
)

func MakeHTTPHandler(ctx context.Context, endpoints Endpoints, logger kitlog.Logger) http.Handler {
	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(errorEncoder),
		kithttp.ServerErrorLogger(logger),
	}
	m := http.NewServeMux()
	m.Handle("/login", kithttp.NewServer(
		ctx,
		endpoints.Login,
		DecodeHTTPLoginRequest,
		EncodeHTTPGenericResponse,
		options...,
	))
	return m
}

func DecodeHTTPLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func EncodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(interface {
		error() error
	}); ok && e.error() != nil {
		errorEncoder(ctx, e.error(), w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	msg := err.Error()

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorWrapper{Error: msg})
}

func errorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

type errorWrapper struct {
	Error string `json:"error"`
}

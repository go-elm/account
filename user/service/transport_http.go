package service

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-elm/account/user"
	"github.com/gorilla/mux"

	"golang.org/x/net/context"
)

var (
	// errBadRoute is used for mux errors
	errBadRoute = errors.New("bad route")
)

func DecodeHTTPNewUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req newUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req.payload); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeHTTPUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r)
	if err != nil {
		return nil, err
	}
	return userRequest{ID: id}, nil
}

func DecodeHTTPUsersRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return usersRequest{}, nil
}

func DecodeHTTPModifyUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r)
	if err != nil {
		return nil, err
	}
	req := modifyUserRequest{ID: id}
	if err := json.NewDecoder(r.Body).Decode(&req.payload); err != nil {
		return nil, err
	}
	return req, nil
}

func idFromRequest(r *http.Request) (user.ID, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return "", errBadRoute
	}
	return user.ID(id), nil
}

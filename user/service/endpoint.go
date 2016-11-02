package service

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/go-elm/account/user"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	NewUser    endpoint.Endpoint
	User       endpoint.Endpoint
	Users      endpoint.Endpoint
	ModifyUser endpoint.Endpoint
}

func MakeNewUserEndpoint(svc user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(newUserRequest)
		u, err := svc.NewUser(ctx, req.payload)
		return newUserResponse{
			User: newResponse(u),
			Err:  err,
		}, nil
	}
}

func MakeUserEndpoint(svc user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(userRequest)
		u, err := svc.User(ctx, req.ID)
		return newUserResponse{
			User: newResponse(u),
			Err:  err,
		}, nil
	}
}

func MakeUsersEndpoint(svc user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		users, err := svc.Users(ctx)
		if err != nil {
			return usersResponse{Err: err}, nil
		}
		var resp usersResponse
		for _, u := range users {
			ur := newResponse(&u)
			resp.Users = append(resp.Users, *ur)
		}
		return resp, nil
	}
}

func MakeModifyUserEndpoint(svc user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(modifyUserRequest)
		u, err := svc.ModifyUser(ctx, req.ID, req.payload)
		return modifyUserResponse{
			User: newResponse(u),
			Err:  err,
		}, nil
	}
}

type newUserRequest struct {
	payload user.Payload
}

type newUserResponse struct {
	User *response `json:"user,omitempty"`
	Err  error     `json:"err,omitempty"`
}

func (r newUserResponse) error() error { return r.Err }

func (r newUserResponse) status() int { return http.StatusCreated }

type userRequest struct {
	ID user.ID
}

type userResponse struct {
	User *response `json:"user,omitempty"`
	Err  error     `json:"err,omitempty"`
}

func (r userResponse) error() error { return r.Err }

type usersRequest struct{}

type usersResponse struct {
	Users []response `json:"err,omitempty"`
	Err   error      `json:"err,omitempty"`
}

func (r usersResponse) error() error { return r.Err }

type modifyUserRequest struct {
	ID      user.ID
	payload user.Payload
}

type modifyUserResponse struct {
	User *response `json:"user,omitempty"`
	Err  error     `json:"err,omitempty"`
}

func (r modifyUserResponse) error() error { return r.Err }

// response has fields used to respond to a client.
type response struct {
	ID        user.ID `json:"user_id"`
	Username  string  `json:"username"`
	FullName  string  `json:"full_name"`
	Email     string  `json:"email"`
	AvatarURL string  `json:"avatar_url,omitempty"`
	Admin     bool    `json:"admin"`
	Enabled   bool    `json:"enabled"`
}

func newResponse(u *user.User) *response {
	return &response{
		ID:        u.ID,
		Username:  u.Username,
		FullName:  u.FullName,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
		Admin:     u.Admin,
		Enabled:   u.Enabled,
	}
}

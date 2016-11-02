package service

import (
	"golang.org/x/net/context"

	"github.com/go-elm/account/auth"
	"github.com/go-elm/account/user"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	Login  endpoint.Endpoint
	Logout endpoint.Endpoint
}

func MakeLoginEndpoint(svc auth.Session) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(loginRequest)
		u, token, err := svc.Login(req.Username, req.Password)
		if err != nil {
			return loginResponse{Err: err}, nil
		}
		return loginResponse{
			User:  newResponse(u),
			Token: &token,
		}, nil
	}
}

type loginRequest struct {
	Username string
	Password string
}

type loginResponse struct {
	User  *response `json:"user,omitempty"`
	Token *string   `json:"token,omitempty"`
	Err   error     `json:"err,omitempty"`
}

func (r loginResponse) error() error { return r.Err }

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

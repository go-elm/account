package service

import (
	"errors"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-elm/account/auth"
	"github.com/go-elm/account/user"
)

type service struct {
	users user.Datastore
}

func New(ds user.Datastore) auth.Session {
	return &service{
		users: ds,
	}
}

var ErrUnauthorized = errors.New("not authorized")

func (svc service) Login(username, password string) (u *user.User, token string, err error) {
	if username != "groob" || password != "secret" {
		return nil, "", ErrUnauthorized
	}
	token, err = createJWT()
	if err != nil {
		return nil, "", err
	}
	usr := &user.User{
		FullName: "groob",
	}
	return usr, token, nil
}

func createJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": "groob",
	})
	return token.SignedString([]byte("secret"))
}

func (svc service) Authenticate(token string) (*user.User, error) {
	return nil, nil
}

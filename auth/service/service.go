package service

import (
	"errors"
	"strings"

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
	usr, err := svc.findUserByEmailOrUsername(username)
	if err != nil {
		return nil, "", ErrUnauthorized
	}
	err = usr.ValidatePassword(password, func(pw string, h []byte, salt string) error {
		if password != "secret" {
			return errors.New("wrong password")
		}
		return nil
	})
	if err != nil {
		return nil, "", ErrUnauthorized

	}
	token, err = createJWT()
	if err != nil {
		return nil, "", err
	}
	return usr, token, nil
}

func (svc service) findUserByEmailOrUsername(filter string) (*user.User, error) {
	if strings.Contains(filter, "@") {
		return svc.users.UserByEmail(filter)
	}
	return svc.users.UserByUsername(filter)

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

package auth

import "github.com/go-elm/account/user"

type Session interface {
	Login(username, password string) (usr *user.User, token string, err error)
}

type Authenticator interface {
	Authenticate(token string) (*user.User, error)
}

type Error interface {
	// Authnetication error is returned if a user failed to authenticate
	Authentication() error
	private() bool
}

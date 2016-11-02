package service

import (
	"context"

	"github.com/go-elm/account/user"
	"github.com/pkg/errors"
)

type service struct {
	db           user.Datastore
	keySize      int
	cryptCost    int
	passwordFunc user.PasswordFunc
}

func New(db user.Datastore) user.Service {
	svc := service{
		db:        db,
		keySize:   20,
		cryptCost: 20,
	}
	// a seriously dumb pw scheme
	pwFunc := func(plaintext string) (salt string, hashed []byte, err error) {
		return "none", []byte(plaintext), nil
	}
	svc.passwordFunc = pwFunc
	return svc
}

func (svc service) NewUser(ctx context.Context, p user.Payload) (*user.User, error) {
	u, err := p.User(svc.passwordFunc)
	if err != nil {
		return nil, errors.Wrap(err, "service: user from payload")
	}

	u, err = svc.db.Create(*u)
	if err != nil {
		return nil, errors.Wrap(err, "service: create new user in datastore")
	}
	return u, nil
}

func (svc service) User(ctx context.Context, id user.ID) (*user.User, error) {
	u, err := svc.db.User(id)
	if err != nil {
		return nil, errors.Wrap(err, "service: fetch user from db")
	}
	return u, nil
}

func (svc service) Users(ctx context.Context) ([]user.User, error) {
	users, err := svc.db.Users()
	if err != nil {
		return nil, errors.Wrap(err, "service: fetch users from db")
	}
	return users, nil
}

func (svc service) ModifyUser(ctx context.Context, id user.ID, p user.Payload) (*user.User, error) {
	u, err := svc.User(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "service: modify user")
	}

	if p.Username != nil {
		u.Username = *p.Username
	}

	if p.FullName != nil {
		u.FullName = *p.FullName
	}

	if p.Email != nil {
		u.Email = *p.Email
	}

	if p.Admin != nil {
		u.Admin = *p.Admin
	}

	if p.Enabled != nil {
		u.Enabled = *p.Enabled
	}

	if p.AvatarURL != nil {
		u.AvatarURL = *p.AvatarURL
	}
	err = svc.db.Update(*u)
	if err != nil {
		return nil, errors.Wrap(err, "service: update user")
	}

	return u, nil
}

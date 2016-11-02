// Package user defines types for managing a website users.
package user

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

// ID is a unique id for a User.
// An ID is typically represented by a UUID.
type ID string

// User is an application user.
type User struct {
	ID        ID
	Username  string
	FullName  string
	Email     string
	Password  []byte
	Salt      string
	AvatarURL string

	CreatedAt time.Time
	UpdatedAt time.Time

	Admin   bool
	Enabled bool
}

// Payload has fields used in a request to create
// or update a user.
type Payload struct {
	Username  *string `json:"username"`
	FullName  *string `json:"full_name"`
	Email     *string `json:"email"`
	AvatarURL *string `json:"avatar_url"`
	Password  *string `json:"password"`
	Admin     *bool   `json:"admin"`
	Enabled   *bool   `json:"enabled"`
}

// User converts a payload into a User.
func (p Payload) User(passwordFn PasswordFunc) (*User, error) {
	u := &User{
		Username: *p.Username,
		FullName: *p.FullName,
		Email:    *p.Email,
		Admin:    p.Admin == nil && *p.Admin,
		Enabled:  true,
	}
	var err error
	u.Salt, u.Password, err = passwordFn(*p.Password)
	if err != nil {
		return nil, errors.Wrap(err, "create user from payload")
	}

	return u, nil
}

// PasswordFunc is a function used to set the user password.
type PasswordFunc func(plaintext string) (salt string, password []byte, err error)

// Service defines methods used to manage users in an application.
type Service interface {
	// NewUser creates a new user from a UserPayload.
	NewUser(ctx context.Context, p Payload) (user *User, err error)

	// User returns a user with a specific UserID.
	User(ctx context.Context, id ID) (user *User, err error)

	// Users returns a list of users.
	Users(ctx context.Context) (users []User, err error)

	// ModifyUser updates a user's properties.
	ModifyUser(ctx context.Context, id ID, p Payload) (user *User, err error)
}

// Datastore manages users in a datastore.
type Datastore interface {
	// CreateUser creates a new user.
	// UserExistsError is returned if a user already exists.
	Create(u User) (user *User, err error)

	// Update updates an existing user.
	// NotFoundError is returned if a user is not found.
	Update(u User) (err error)

	// User retrieves a single user from the datastore.
	User(id ID) (user *User, err error)

	// Users retrieves a list of users from the datastore.
	Users() (users []User, err error)
}

// Error encapsulates the possible errors returned by
// the Service and Datastore.
type Error interface {
	// Exits error is returned if a user already exists and cannot be created.
	Exists() error

	// NotFound error is returned if a specific user cannot be found.
	NotFound() error

	// private restricts this interface so that it cannot be implemented.
	// Error only exists for documentation.
	private() bool
}

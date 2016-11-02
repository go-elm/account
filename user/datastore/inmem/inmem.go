package inmem

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/go-elm/account/user"
	"github.com/pkg/errors"
)

type inmem struct {
	mtx   sync.RWMutex
	users map[user.ID]*user.User
}

func New() *inmem {
	return &inmem{
		users: make(map[user.ID]*user.User),
	}
}

func (db *inmem) Create(u user.User) (*user.User, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	uuid, err := uuidGen()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create uuid")
	}
	u.ID = uuid
	db.users[u.ID] = &u

	return &u, nil
}

func (db *inmem) Updateuser(u user.User) error {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	if ok, err := db.exists(u); !ok {
		return err
	}
	db.users[u.ID] = &u
	return nil
}

func (db *inmem) User(id user.ID) (*user.User, error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()

	u, ok := db.users[id]
	if !ok {
		return nil, &notFoundError{
			u.ID, u.Username,
		}
	}
	return u, nil
}

func (db *inmem) Users() ([]user.User, error) {
	var unsorted []user.User
	for _, user := range db.users {
		unsorted = append(unsorted, *user)
	}
	return unsorted, nil
}

func (db *inmem) exists(u user.User) (bool, error) {
	if _, ok := db.users[u.ID]; !ok {
		return false, &notFoundError{
			u.ID, u.Username,
		}
	}
	return true, nil
}

type notFoundError struct {
	ID       user.ID
	Username string
}

func (e *notFoundError) Error() string {
	return fmt.Sprintf("user %s, with id %s not found in inmem datastore.", e.Username, e.ID)
}

func (e *notFoundError) NotFound() error {
	return e
}

func uuidGen() (user.ID, error) {
	unix := uint32(time.Now().UTC().Unix())

	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}

	return user.ID(fmt.Sprintf("%08x-%04x-%04x-%04x-%04x%08x",
			unix, b[0:2], b[2:4], b[4:6], b[6:8], b[8:])),
		nil
}

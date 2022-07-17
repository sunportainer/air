package air

import "errors"

type (
	// DataStore defines the interface to manage the data
	DataStore interface {
		Open() (newStore bool, err error)
		Init() error
		Close() error
		Settings() DBService
	}

	DBService interface {
		UserByEmail(email string) (*User, error)
		CreateUser(user *User) error
	}
	RedisService interface {
		UserByEmail(name string) (*User, error)
		CreateUser(user *User) error
	}
	User struct {
		ID        int    `json:"ID"`
		Email     string `json:"Email"`
		Password  string `json:"Password"`
		FirstName string `json:"FirstName"`
		LastName  string `json:"LastName"`
	}

	MockDB struct {
		Users []User `json:"users"`
	}
)

var (
	ErrObjectNotFound = errors.New("object not found inside the database")
)

const KeyServerAddr = "serverAddr"

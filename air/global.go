package air

import "errors"

type (
	User struct {
		ID        int    `json:"ID"`
		Email     string `json:"Email"`
		Login     string `json:"login,omitempty"`
		Password  string `json:"Password,omitempty"`
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

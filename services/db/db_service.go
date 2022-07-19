package db

import (
	"com.nzair.user/air"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type DBService interface {
	UserByEmail(email string) (*air.User, error)
	CreateUser(user *air.User) error
}

//here we use json to mock one postsqlgre database
type Service struct {
	db string
}

func NewDBService(dbFile string) (DBService, error) {
	service := Service{dbFile}
	if !dbExists(dbFile) {
		file := `{"users": []}`
		ioutil.WriteFile(dbFile, []byte(file), os.FileMode(0755))
	}
	return &service, nil
}

//query user info from database
func (service *Service) UserByEmail(email string) (*air.User, error) {
	db, err := service.loadDB()
	if err != nil {
		log.Println("Failed to load DB file")
		return nil, err
	}
	for _, user := range db.Users {
		if strings.EqualFold(email, user.Email) {
			return &user, nil
		}
	}
	return nil, air.ErrObjectNotFound

}

//create one user, will use it in later checkin
func (service *Service) CreateUser(user *air.User) error {
	db, err := service.loadDB()
	if err != nil {
		log.Println("Failed to load DB file")
		return err
	}
	db.Users = append(db.Users, *user)
	err = service.flushDB(db)
	if err != nil {
		log.Println("Failed to create user in DB")
		return err
	}
	return nil
}

//load json-based db file to struct
func (service *Service) loadDB() (*air.MockDB, error) {
	filePtr, err := os.Open(service.db)
	if err != nil {
		log.Println("Failed to read DB file")
		return nil, err
	}
	defer filePtr.Close()
	var db air.MockDB
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&db)
	if err != nil {
		log.Println("Failed to decode db file", err.Error())
	}
	return &db, nil
}

//write structs back to json-based db file
func (service *Service) flushDB(db *air.MockDB) error {

	filePtr, err := os.OpenFile(service.db, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		log.Println("Failed to read DB file")
		return err
	}
	defer filePtr.Close()
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(db)
	if err != nil {
		fmt.Println("Encode error when flush database", err.Error())
	}
	return nil
}

func dbExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

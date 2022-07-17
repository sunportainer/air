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

//here we use json 代替数据库
type Service struct {
}

func NewDBService() (*Service, error) {
	service := &Service{}
	if !dbExists("./db.json") {
		file := `{"users": []}`
		ioutil.WriteFile("./db.json", []byte(file), os.FileMode(0755))
	}
	return service, nil

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
	filePtr, err := os.Open("./db.json")
	if err != nil {
		log.Println("Failed to read DB file")
		return nil, err
	}
	defer filePtr.Close()
	var db air.MockDB
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&db)
	if err != nil {
		log.Println("Failed to decode db file", err.Error())
	}
	return &db, nil
}

//write structs back to json-based db file
func (service *Service) flushDB(db *air.MockDB) error {
	filePtr, err := os.Open("./db.json")
	if err != nil {
		log.Println("Failed to read DB file")
		return err
	}
	defer filePtr.Close()
	// 创建Json编码器
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(db)
	if err != nil {
		fmt.Println("编码错误", err.Error())
	} else {
		fmt.Println("编码成功")
	}
	return nil
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

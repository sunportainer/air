package handler

import (
	"com.nzair.user/air"
	"com.nzair.user/services/db"
	"com.nzair.user/services/myredis"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"io"
	"log"
	"net/http"
	"strings"
)

// UserHandler
type UserHandler struct {
	DBService    db.DBService
	RedisService myredis.RedisService
}

//create new handler for http req
func NewHandler() (*UserHandler, error) {
	dbService, err := db.NewDBService("./db.json")
	if err != nil {
		log.Println("Failed to create handler")
		return nil, err
	}
	redisService, err := myredis.NewRedisService()
	if err != nil {
		log.Println("Failed to create handler")
		return nil, err
	}
	handler := &UserHandler{DBService: dbService, RedisService: redisService}
	return handler, nil

}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		h.handlerGet(w, req)
	} else {
		h.handlerPost(w, req)
	}
}

func (h *UserHandler) handlerGet(w http.ResponseWriter, req *http.Request) {
	reqEmail := req.URL.Query().Get("email")
	reqPassword := req.URL.Query().Get("pwd")
	if reqPassword == "" || reqEmail == "" {
		log.Println("Request field missing")
		w.Header().Set("x-missing-field", "pwd or username")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "BadRequest: request field missing.\n")
		return
	}
	//query user detail from myredis
	user, err := h.RedisService.UserByEmail(reqEmail)
	if err != nil {
		//if not found then try to hit db directly
		if err == redis.Nil {
			user, err = h.DBService.UserByEmail(reqEmail)
			if err != nil {
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, "User not found.\n")
				return
			}
			//get user in db, then set to myredis and return
			err = h.RedisService.CreateUser(user)
			if err != nil {
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, "Failed to refresh myredis.\n")
				return
			}
		}
	}
	//validate password
	if strings.Compare(user.Password, reqPassword) == 0 {
		w.WriteHeader(http.StatusOK)
		user.Password = ""
		user.Login = ""
		h.JSON(w, user)
	} else {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Failed to login.\n")
	}
}

func (h *UserHandler) handlerPost(w http.ResponseWriter, req *http.Request) {
	log.Printf("The request is from %s", req.RemoteAddr)
	//To make it simple, just ignore validate here
	firstName := req.PostFormValue("firstName")
	lastName := req.PostFormValue("lastName")
	email := req.PostFormValue("email")
	password := req.PostFormValue("pwd")
	login := req.RemoteAddr
	newUser := &air.User{ID: 1, FirstName: firstName, LastName: lastName, Password: password, Login: login, Email: email}
	err := h.DBService.CreateUser(newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Failed to create user.\n")
	}
}

//write JSON to response
func (h *UserHandler) JSON(rw http.ResponseWriter, data interface{}) error {
	rw.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(rw).Encode(data)
	if err != nil {
		return errors.New("Unable to write JSON response")
	}
	return nil
}

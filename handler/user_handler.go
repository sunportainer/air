package handler

import (
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
	dbService, err := db.NewDBService()
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

	reqEmail := req.URL.Query().Get("email")
	reqPassword := req.URL.Query().Get("pwd")
	if reqPassword == "" || reqEmail == "" {
		log.Println("Request field missing")
		w.Header().Set("x-missing-field", "pwd or username")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Request field missing.\n")
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
		h.JSON(w, user)
	} else {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Failed to login.\n")
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

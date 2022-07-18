package handler

import (
	"com.nzair.user/air"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type MockDBService struct {
	handleQuery  func() (*air.User, error)
	handleCreate func() error
}
type MockRedisService struct {
	handleQuery  func() (*air.User, error)
	handleCreate func() error
}

func (res *MockRedisService) UserByEmail(email string) (*air.User, error) {
	return res.handleQuery()
}

func (res *MockRedisService) CreateUser(user *air.User) error {
	return res.handleCreate()
}
func (res *MockDBService) UserByEmail(email string) (*air.User, error) {
	return res.handleQuery()
}

func (res *MockDBService) CreateUser(user *air.User) error {
	return res.handleCreate()
}

func TestUserHandler_JSON(t *testing.T) {
	//TBA
}

func TestUserHandler_ServeHTTP_BadRequest(t *testing.T) {
	handler := &UserHandler{
		&MockDBService{},
		&MockRedisService{},
	}
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Response should be 400, but get %v", resp.StatusCode)
	}
}

func TestUserHandler_ServeHTTP_FromDB(t *testing.T) {
	mockDB := &MockDBService{}
	mockRedis := &MockRedisService{}
	handler := &UserHandler{
		mockDB,
		mockRedis,
	}
	hashPasswdBytes := md5.Sum([]byte("P@ssw0rd"))
	hashPasswd := string(hashPasswdBytes[:])
	params := url.Values{}
	params.Add("email", "tester@gmail.com")
	params.Add("pwd", hashPasswd)

	mockDB.handleQuery = func() (*air.User, error) {
		return &air.User{ID: 1, FirstName: "tester", LastName: "New Zealand", Email: "tester@gmail.com", Password: hashPasswd}, nil
	}

	mockRedis.handleQuery = func() (*air.User, error) {
		return nil, redis.Nil
	}
	mockRedis.handleCreate = func() error {
		return nil
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:3333?%s", params.Encode()), nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
	user := &air.User{}
	json.Unmarshal(w.Body.Bytes(), user)
	if user.ID != 1 {
		t.Errorf("User ID should equal to 1, but get %v", user.ID)
	}
	if strings.Compare(user.Email, "tester@gmail.com") != 0 {
		t.Errorf("Emails are different, %s != %s", user.Email, "tester@gmail.com")
	}
}

func TestUserHandler_ServeHTTP_FromRedis(t *testing.T) {
	mockDB := &MockDBService{}
	mockRedis := &MockRedisService{}
	handler := &UserHandler{
		mockDB,
		mockRedis,
	}
	hashPasswdBytes := md5.Sum([]byte("P@ssw0rd"))
	hashPasswd := string(hashPasswdBytes[:])
	params := url.Values{}
	params.Add("email", "tester@gmail.com")
	params.Add("pwd", hashPasswd)

	mockRedis.handleQuery = func() (*air.User, error) {
		return &air.User{ID: 1, FirstName: "tester", LastName: "New Zealand", Email: "tester@gmail.com", Password: hashPasswd}, nil
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:3333?%s", params.Encode()), nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
	user := &air.User{}
	json.Unmarshal(w.Body.Bytes(), user)
	if user.ID != 1 {
		t.Errorf("User ID should equal to 1, but get %v", user.ID)
	}
	if strings.Compare(user.Email, "tester@gmail.com") != 0 {
		t.Errorf("Emails are different, %s != %s", user.Email, "tester@gmail.com")
	}
}

package myredis

import (
	"com.nzair.user/air"
	"encoding/json"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"strings"
	"testing"
)

func TestService_UserByEmail(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()
	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	key := "tester@gmail.com"
	tester := &air.User{ID: 1, FirstName: "tester", LastName: "New Zealand", Email: "tester@gmail.com"}
	testData, _ := json.Marshal(tester)
	s.Set(key, string(testData))
	service := &Service{
		redisClient: rdb,
	}

	user, _ := service.UserByEmail(key)
	if strings.Compare(user.Email, key) != 0 {
		t.Error("Email address do not match")
	}
}

func TestService_CreateUser(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()
	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(), // mock redis server的地址
	})
	service := &Service{
		redisClient: rdb,
	}
	tester := &air.User{ID: 1, FirstName: "tester", LastName: "New Zealand", Email: "tester@gmail.com"}
	err = service.CreateUser(tester)
	if err != nil {
		t.Error("Failed to put user into redis")
	}
}

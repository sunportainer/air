package myredis

import (
	"com.nzair.user/air"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type (
	RedisService interface {
		CreateUser(user *air.User) error
		UserByEmail(email string) (*air.User, error)
	}
	Service struct {
		redisClient *redis.Client
	}
)

// init myredis connection to single node mode
func NewRedisService() (RedisService, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: 10, // size of pool
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Println("Redis is not working")
		return nil, err
	}
	newService := &Service{redisClient: redisClient}
	return newService, nil
}

//query user from redis
func (res *Service) UserByEmail(email string) (*air.User, error) {
	ctx := context.Background()

	val, err := res.redisClient.Get(ctx, email).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Println("keys does not exist")
		}
		return nil, err
	}
	var user air.User
	err = json.Unmarshal([]byte(val), &user)
	if nil != err {
		log.Println("Failed to parse result")
		return nil, err
	}
	return &user, nil
}

//Set user to redis
func (res *Service) CreateUser(user *air.User) error {
	ctx := context.Background()
	data, _ := json.Marshal(user)

	err := res.redisClient.Set(ctx, user.Email, data, 0).Err()
	if err != nil {
		log.Println("Failed to create user in myredis")
		return err
	}
	return nil
}

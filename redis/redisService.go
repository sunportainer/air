package redis

import (
	"com.nzair.user/air"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type Service struct {
	redisClient *redis.Client
}

// init redis connection to single node mode
func (res *Service) initClient() (err error) {
	res.redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: 10, // size of pool
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = res.redisClient.Ping(ctx).Result()
	return err
}

func (res *Service) Test() {
	ctx := context.Background()
	user := air.User{
		ID:        1,
		Email:     "eee@gmail.com",
		FirstName: "ERic",
		LastName:  "Sun",
	}
	//json序列化
	datas, _ := json.Marshal(user)
	//rebytes, _ := redis.Bytes(conn.Do("get", "struct3"))
	//json反序列化
	//object := &TestStruct{}
	//json.Unmarshal(rebytes, object)

	err := res.redisClient.Set(ctx, user.Email, datas, 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := res.redisClient.Get(ctx, user.Email).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			fmt.Println("keys does not exist")
		} else {
			panic(err)
		}
	}
	fmt.Println("key", val)
	val2, err := res.redisClient.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
}

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
	//反序列化,这里用指针行吗？
	err = json.Unmarshal([]byte(val), &user)
	if nil != err {
		log.Println("Failed to parse result")
		return nil, err
	}
	return &user, nil
}

func (res *Service) CreateUser(user *air.User) error {
	ctx := context.Background()
	//json序列化
	data, _ := json.Marshal(user)
	err := res.redisClient.Set(ctx, user.Email, data, 0).Err()
	if err != nil {
		log.Println("Failed to create user in redis")
		return err
	}
	return nil
}

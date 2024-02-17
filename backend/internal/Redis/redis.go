package Redis

import (
	"github.com/go-redis/redis"
	"log"
)

var Client *redis.Client

func InitRedis(addr, password string) {
	log.Println("Connected to Database...")
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	log.Println("Success")
}

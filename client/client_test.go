package client

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestNewClient(t *testing.T) {
	c, err := New("localhost:5001")

	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	for i := range 10 {
		if err := c.Set(context.TODO(), fmt.Sprintf("foo_%d", i), fmt.Sprintf("bar_%d", i)); err != nil {
			log.Fatal(err)
		}
		val, err := c.Get(context.TODO(), fmt.Sprintf("foo_%d", i))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("go back value", val)
	}
}

func TestNewClientRedisClient(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:5001",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Println(rdb)
	fmt.Print("dddddd")
	err := rdb.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	// val, err := rdb.Get(context.TODO(), "key").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("key", val)

}

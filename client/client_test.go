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

	if err := c.Set(context.TODO(), "foo", 69); err != nil {
		log.Fatal(err)
	}
	val, err := c.Get(context.TODO(), "foo")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("go back value", val)

}

func TestOfficailRedisClient(t *testing.T) {
	go func ()  {
		
	}()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:5001",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Println(rdb)
	fmt.Print("dddddd")
	err := rdb.Set(context.Background(), "foo", "bar", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(context.TODO(), "foo").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("get value back", val)

}

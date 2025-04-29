package client

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
)

func TestNewClient(t *testing.T) {
	c, err := New("localhost:5001")

	if err != nil {
		log.Fatal(err)
	}

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

func TestNewClients(t *testing.T) {
	nClients := 10
	wg := sync.WaitGroup{}
	wg.Add(nClients)
	for i := 0; i < nClients; i++ {
		go func(it int) {
			c, err := New("localhost:5001")

			if err != nil {
				log.Fatal(err)
			}
			defer c.Close()
			key := fmt.Sprintf("client_%d", it)
			value := fmt.Sprintf("client_%d", it)

			if err := c.Set(context.TODO(), key, value); err != nil {
				log.Fatal(err)
			}
			val, err := c.Get(context.TODO(), key)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("go back value form client %s", val)
			wg.Done()
		}(i)

	}
	wg.Wait()

}

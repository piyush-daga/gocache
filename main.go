package main

import (
	"fmt"
	"go-cache/gocache"
	"time"
)

func main() {
	// c := gocache.NewCacheWithSyncMap()
	// c := gocache.NewCacheWithMutex(2)
	c := gocache.NewCacheWithEvictionPolicy(2)

	c.Set("ok", "testing 123", time.Second*2)

	go func() {
		for {
			_, _ = c.Get("ok")
			c.Set("ok2", "testing 456", time.Second*3)
		}
	}()

	// Simulating race-conditions
	go func() {
		for {
			c.Set("ok1", "testing 789", time.Second*2)
		}
	}()

	time.Sleep(5 * time.Second)
	// As a test
	c.Get("ok")

	// fmt.Println(c.List())

	fmt.Println(c.Get("ok"))
	fmt.Println(c.Get("ok1"))
	fmt.Println(c.Get("ok2"))
}

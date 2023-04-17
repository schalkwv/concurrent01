package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	// Start 5 goroutines.
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(goroutineNum int) {
			defer wg.Done()
			// Wait for a user ID to be passed and print it out.
			for userId := range userIds {
				sleepTime := time.Duration(rand.Intn(7)) * time.Second
				fmt.Printf("routine %v User ID: %v sleeping %v\n", goroutineNum, userId, sleepTime)
				time.Sleep(sleepTime)
			}
		}(i)
	}

	// Loop through 100 user IDs and pass them to the goroutines only if there is one available.
	for i := 1; i <= 100; i++ {
		userIds <- fmt.Sprintf("user%d", i)
	}

	close(userIds)
	wg.Wait()
}

var userIds = make(chan string, 1)

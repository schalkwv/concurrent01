package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type application struct {
	hub *Hub
}

func (app *application) Broadcast(goroutineNum string, userID string) {
	data := map[string]string{
		"goroutine": goroutineNum,
		"userID":    userID,
	}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		//fmt.Println(string(jsonStr))
		app.hub.broadcast <- message{
			roomID: "1",
			data:   jsonStr,
		}
	}
}

func (app *application) handleWebSocket(c echo.Context) error {
	serveWs(app.hub, c.Response(), c.Request(), "1")
	return nil
}
func main() {
	app := &application{
		hub: newHub(),
	}
	go app.hub.run()
	// Start Echo server and WebSocket endpoint.
	e := echo.New()
	e.Use(middleware.CORS())
	e.GET("/ws", app.handleWebSocket)
	go e.Start(":8085")

	// Start 5 goroutines.
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(goroutineNum int) {
			defer wg.Done()
			// Wait for a user ID to be passed, sleep for a random number of seconds, and then print it out.
			for userId := range userIds {
				sleepTime := time.Duration(rand.Intn(5)) * time.Second
				fmt.Printf("++start routine %v User ID: %v sleeping %v\n", goroutineNum, userId, sleepTime)
				time.Sleep(sleepTime)
				fmt.Printf("--end routine %v User ID: %v sleeping %v\n", goroutineNum, userId, sleepTime)
				app.Broadcast(strconv.Itoa(goroutineNum), userId)

			}
		}(i)
	}

	// Loop through 100 user IDs and pass them to the goroutines only if there is one available.
	for i := 1; i <= 50; i++ {
		userIds <- fmt.Sprintf("%d", i)
	}

	close(userIds)
	wg.Wait()
	//give the socket time to broadcast
	time.Sleep(5 * time.Second)
}

var userIds = make(chan string, 1)

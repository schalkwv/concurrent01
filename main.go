package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	go cleanUpWebSocketClients()
	// Start Echo server and WebSocket endpoint.
	e := echo.New()
	e.GET("/ws", handleWebSocket)
	go e.Start(":8085")

	// Start 5 goroutines.
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(goroutineNum int) {
			defer wg.Done()
			// Wait for a user ID to be passed, sleep for a random number of seconds, and then print it out.
			for userId := range userIds {
				sleepTime := time.Duration(rand.Intn(5)) * time.Second // Sleep for up to 5 seconds.
				time.Sleep(sleepTime)
				//message := fmt.Sprintf("Goroutine %v: User ID: %v", goroutineNum, userId)
				//broadcastWebSocket(message)
				fmt.Printf("routine %v User ID: %v sleeping %v\n", goroutineNum, userId, sleepTime)
				broadcast(strconv.Itoa(goroutineNum), userId)

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

func handleWebSocket(c echo.Context) error {
	// Upgrade WebSocket connection.
	fmt.Printf("handleWebSocket")
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	websocketClientsLock.Lock()
	websocketClients[ws] = true
	websocketClientsLock.Unlock()

	// Continuously read messages from WebSocket and discard them.
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}

	return nil
}

func broadcast(goroutineNum string, userID string) {
	message := map[string]string{
		"goroutine": goroutineNum,
		"userID":    userID,
	}
	for conn := range websocketClients {
		err := conn.WriteJSON(message)
		if err != nil {
			log.Printf("error broadcasting message: %v", err)
			conn.Close()
			delete(websocketClients, conn)
		}
	}
}

var websocketClients = make(map[*websocket.Conn]bool)
var websocketClientsLock sync.Mutex

func init() {
	// Set random seed for sleep times.
	rand.Seed(time.Now().UnixNano())
}

//func init() {
//	// Periodically clean up closed WebSocket connections.
//	go func() {
//		for {
//			time.Sleep(5 * time.Second)
//			websocketClientsLock.Lock()
//			for client := range websocketClients {
//				if !client.IsAlive() {
//					client.Close()
//					delete(websocketClients, client)
//				} else {
//					client.SetReadDeadline(time.Now().Add(10 * time.Second))
//					client.SetWriteDeadline(time.Now().Add(10 * time.Second))
//				}
//			}
//			websocketClientsLock.Unlock()
//		}
//	}()
//}

func cleanUpWebSocketClients() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		websocketClientsLock.Lock()
		for conn := range websocketClients {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Printf("error reading message from WebSocket client: %v", err)
				conn.Close()
				delete(websocketClients, conn)
			}
		}
		websocketClientsLock.Unlock()
	}
}

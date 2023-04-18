package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		//origin := r.Header.Get("Origin")
		return true //origin == "http://127.0.0.1:8080"
	},
}

type application struct {
	hub *Hub
}

//func broadcast(goroutineNum string, userID string) {
//	message := map[string]string{
//		"goroutine": goroutineNum,
//		"userID":    userID,
//	}
//	for conn := range websocketClients {
//		err := conn.WriteJSON(message)
//		if err != nil {
//			log.Printf("error broadcasting message: %v", err)
//			conn.Close()
//			delete(websocketClients, conn)
//		}
//	}
//}

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
				sleepTime := time.Duration(rand.Intn(5)) * time.Second // Sleep for up to 5 seconds.
				//sleepTime := time.Second * 3
				time.Sleep(sleepTime)
				//message := fmt.Sprintf("Goroutine %v: User ID: %v", goroutineNum, userId)
				//broadcastWebSocket(message)
				//fmt.Printf("routine %v User ID: %v sleeping %v\n", goroutineNum, userId, sleepTime)
				//broadcast(strconv.Itoa(goroutineNum), userId)
				app.Broadcast(strconv.Itoa(goroutineNum), userId)

			}
		}(i)
	}

	// Loop through 100 user IDs and pass them to the goroutines only if there is one available.
	for i := 1; i <= 200; i++ {
		userIds <- fmt.Sprintf("%d", i)
	}

	close(userIds)
	wg.Wait()
}

var userIds = make(chan string, 1)

//func handleWebSocket(c echo.Context) error {
//	// Upgrade WebSocket connection.
//	fmt.Printf("handleWebSocket")
//	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
//	if err != nil {
//		fmt.Println("error in ws handler")
//		fmt.Println(err)
//		return err
//	}
//	defer ws.Close()
//
//	websocketClientsLock.Lock()
//	websocketClients[ws] = true
//	websocketClientsLock.Unlock()
//
//	// Continuously read messages from WebSocket and discard them.
//	for {
//		_, _, err := ws.ReadMessage()
//		if err != nil {
//			break
//		}
//	}
//
//	return nil
//}

//func handleWebSocket(c echo.Context) error {
//	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
//	if err != nil {
//		return err
//	}
//	defer conn.Close()
//
//	// create channel to queue messages for this connection
//	messageChan := make(chan []byte, 100)
//
//	// start goroutine to write messages to the connection
//	go func() {
//		for message := range messageChan {
//			err := conn.WriteMessage(websocket.TextMessage, message)
//			if err != nil {
//				log.Printf("error writing message to WebSocket connection: %v", err)
//				break
//			}
//		}
//	}()
//
//	// add connection to map of active WebSocket clients
//	websocketClientsLock.Lock()
//	websocketClients[conn] = messageChan
//	websocketClientsLock.Unlock()
//
//	// read from WebSocket connection and broadcast to all clients
//	for {
//		_, _, err := conn.ReadMessage()
//		if err != nil {
//			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
//				log.Printf("error reading message from WebSocket connection: %v", err)
//			}
//			break
//		}
//		// ...
//	}
//
//	// remove connection from map of active WebSocket clients
//	websocketClientsLock.Lock()
//	delete(websocketClients, conn)
//	close(messageChan)
//	websocketClientsLock.Unlock()
//
//	return nil
//}

func init() {
	// Set random seed for sleep times.
	rand.Seed(time.Now().UnixNano())
}

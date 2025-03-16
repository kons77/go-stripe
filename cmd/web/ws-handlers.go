package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketConnection struct {
	*websocket.Conn
}

// WsPayload defines the structure of the data received from clients
type WsPayload struct {
	Action      string              `json:"action"`
	Message     string              `json:"message"`
	UserName    string              `json:"username"`
	MessageType string              `json:"message_type"`
	UserID      int                 `json:"user_id"`
	Conn        WebSocketConnection `json:"-"`
}

// WsJsonResponse defines the structure of the response sent to clients
type WsJsonResponse struct {
	Action  string `json:"action"`
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
}

// clients stores active WebSocket connections mapped to usernames
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // how you secure your socket connection
}

var clients = make(map[WebSocketConnection]string)

// wsChan is a channel for handling incoming WebSocket messages
var wsChan = make(chan WsPayload)

// WsEndPoint handles new WebSocket connection requests
func (app *application) WsEndPoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	//defer ws.Close()

	app.infoLog.Println("Client connected from" + r.RemoteAddr)

	var response WsJsonResponse
	response.Message = "Connected to server"

	err = ws.WriteJSON(response)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	go app.ListenForWS(&conn)
}

// ListenForWS listens for incoming WebSocket messages from a specific connection
func (app *application) ListenForWS(conn *WebSocketConnection) {
	// Recover from panics to prevent the application from crashing
	defer func() {
		if r := recover(); r != nil {
			app.errorLog.Println("ERROR:", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	// infinite loop to continuously listen for messages from the client
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothing
			break
		} else {
			payload.Conn = *conn // Attach the connection info to the payload
			wsChan <- payload    // Send the payload to the channel for processing
		}
	}
}

// ListenToWsChannel continuously listens for messages from wsChan and processes them
func (app *application) ListenToWsChannel() {
	var response WsJsonResponse
	for {
		e := <-wsChan
		switch e.Action {
		case "deleteUser":
			// Handle user deletion by notifying all clients
			response.Action = "logout"
			response.Message = "Your account has been deleted"
			response.UserID = e.UserID
			app.broadcastToAll(response)

		default:
		}
	}
}

// broadcastToAll sends a message to all connected WebSocket clients
func (app *application) broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		// broadcast to every connected client
		err := client.WriteJSON(response)
		if err != nil {
			app.errorLog.Printf("WebSocket error on %s: %s", response.Action, err)
			_ = client.Close()
			delete(clients, client)
		}
	}
}

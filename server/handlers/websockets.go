package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/info344-a17/typy-bird/server/models"
)

const maxPlayers = 4

//WebSocketsHandler is a handler for WebSocket upgrade requests
type WebSocketsHandler struct {
	notifier *Notifier
	gameRoom *models.GameRoom
	upgrader *websocket.Upgrader
}

//NewWebSocketsHandler constructs a new WebSocketsHandler
func NewWebSocketsHandler(notifer *Notifier, gameroom *models.GameRoom) *WebSocketsHandler {
	return &WebSocketsHandler{
		notifier: notifer,
		gameRoom: gameroom,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

//ServeHTTP implements the http.Handler interface for the WebSocketsHandler
func (wsh *WebSocketsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("error upgrading websocket: %v", err), http.StatusInternalServerError)
		return
	}
	if len(wsh.notifier.clients) != maxPlayers {
		go wsh.notifier.AddClient(conn)
	} else {
		http.Error(w, "game room is full", http.StatusConflict)
	}
}

//Notifier is an object that handles WebSocket notifications
type Notifier struct {
	clients  []*websocket.Conn
	eventQ   chan []byte
	mx       *sync.RWMutex
	gameRoom *models.GameRoom
}

//NewNotifier constructs a new Notifier
func NewNotifier(gameroom *models.GameRoom) *Notifier {
	notifier := &Notifier{
		clients:  make([]*websocket.Conn, 0),
		eventQ:   make(chan []byte),
		mx:       &sync.RWMutex{},
		gameRoom: gameroom,
	}
	go notifier.start()
	return notifier
}

//AddClient adds a new client to the Notifier
func (n *Notifier) AddClient(client *websocket.Conn) {
	n.mx.Lock()
	n.clients = append(n.clients, client)
	n.mx.Unlock()
	for {
		if _, _, err := client.NextReader(); err != nil {
			n.removeClient(client)
			break
		}
	}
}

//Notify broadcasts the event to all WebSocket clients
func (n *Notifier) Notify(event []byte) {
	log.Printf("adding event to the queue")
	n.eventQ <- event
}

//start starts the notification loop
func (n *Notifier) start() {
	log.Println("starting notifier loop")
	for {
		event := <-n.eventQ
		for _, client := range n.clients {
			err := client.WriteMessage(websocket.TextMessage, event)
			if err != nil {
				n.removeClient(client)
				fmt.Errorf("error writing message: %v", err)
				return
			}
		}
	}
}

func (n *Notifier) removeClient(client *websocket.Conn) {
	client.Close()
	for i, c := range n.clients {
		if client == c {
			n.mx.Lock()
			n.clients = append(n.clients[:i], n.clients[i+1:]...)
			n.gameRoom.Players = append(n.gameRoom.Players[:i], n.clients[i+1]...)
			n.mx.Unlock()
		}
	}
}

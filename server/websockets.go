package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

//WebSocketsHandler is a handler for WebSocket upgrade requests
type WebSocketsHandler struct {
	notifier *Notifier
	upgrader *websocket.Upgrader
}

//NewWebSocketsHandler constructs a new WebSocketsHandler
func NewWebSocketsHandler(notifer *Notifier) *WebSocketsHandler {
	return &WebSocketsHandler{
		notifier: notifer,
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
		http.Error(w, fmt.Sprintf("error upgrading websocker: %v", err), http.StatusInternalServerError)
		return
	}
	wsh.notifier.AddClient(conn)
}

//Notifier is an object that handles WebSocket notifications
type Notifier struct {
	clients []*websocket.Conn
	eventQ  chan []byte
	mx      *sync.RWMutex
}

//NewNotifier constructs a new Notifier
func NewNotifier() *Notifier {
	notifier := &Notifier{
		clients: make([]*websocket.Conn, 0),
		eventQ:  make(chan []byte),
		mx:      &sync.RWMutex{},
	}
	go notifier.start()
	return notifier
}

//AddClient adds a new client to the Notifier
func (n *Notifier) AddClient(client *websocket.Conn) {
	n.mx.Lock()
	n.clients = append(n.clients, client)
	fmt.Println(n.clients)
	n.mx.Unlock()
	for {
		if _, _, err := client.NextReader(); err != nil {
			client.Close()
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
				fmt.Errorf("error writing message: %v", err)
				return
			}
		}
	}
}

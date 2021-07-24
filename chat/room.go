package main

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize}

type room struct {
	// forward is the channel that holds incoming messages
	// that will be forwarded to clients
	forward chan []byte
	// join is a channel for clients wanting to join the room
	join chan *client
	// leave is a channel for clients wanting to leave the room
	leave chan *client
	// clients holds all the current clients in this room
	clients map[*client]bool
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// joining
			r.clients[client] = true
		case client := <-r.leave:
			// leave room
			delete(r.clients, client)
			// close the channel for the client that left
			close(client.send)
		case msg := <-r.forward:
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	// ask room to join newly connected client
	r.join <- client
	// make sure to disconnect the client when it leaves
	defer func() { r.leave <- client }()
	// start client write thread
	go client.write()
	// start client read
	client.read()
}

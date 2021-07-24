package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	// socket is websocket for this client
	socket *websocket.Conn
	// channel where messages are sent
	send chan []byte
	// room client is chatting in
	room *room
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		// read from websocket
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		// send msg to room
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	// read from the send channel and write it out to the websocket
	for msg := range c.send {
		// send message
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}

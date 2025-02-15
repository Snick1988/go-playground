package main

import (
    "fmt"
    "github.com/gorilla/websocket"
)

// client represents a single chatting user.
type client struct {
    // socket is the web socket for this client.
    socket *websocket.Conn
    // send is a channel on which messages are sent.
    send chan []byte
    // room is the room this client is chatting in.
    room *room
    // client name
    name string
}

func (c *client) read() {
    for {
        if _, msg, err := c.socket.ReadMessage(); err == nil {
            c.room.forward <- []byte(fmt.Sprintf("%s: %s", c.name, msg))
        } else {
            break
        }
    }
    c.socket.Close()
}

func (c * client) write() {
    for msg := range c.send {
        if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
            break
        }
    }
    c.socket.Close()
}

package main

import (
    "log"
    "fmt"
    "math/rand"
    "net/http"
    "github.com/gorilla/websocket"
)

type room struct {
    // forward is a channel that holds incoming messages
    // that should be forwarded to the other clients.
    forward chan []byte
    // Join channel
    join chan *client
    // Leave channel
    leave chan *client
    // List of current clients
    clients map[*client]bool
}

const (
    socketBufferSize  = 1024
    messageBufferSize = 256
)

// newRoom makes a new room that is ready to go.
func newRoom() *room {
    return &room{
        forward: make(chan []byte, messageBufferSize),
        join:    make(chan *client),
        leave:   make(chan *client),
        clients: make(map[*client]bool),
    }
}

func (r *room) run() {
    for {
        select {
        case client := <- r.join:
            // Client joins
            r.clients[client] = true
            r.forward <- []byte(fmt.Sprintf("%s joined the room", client.name))
        case client := <- r.leave:
            // Client leaves
            delete(r.clients, client)
            close(client.send)
        case msg := <- r.forward:
            log.Println("Received a message:", string(msg))
            for client := range r.clients {
                select {
                case client.send <- msg:
                // Send message
                default:
                    delete(r.clients, client)
                    close(client.send)
                }
            }
        }
    }
}

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    socket, err := upgrader.Upgrade(w, req, nil)
    if err != nil {
        log.Fatal("ServeHTTP:", err)
        return
    }

    names := make([]string, 0)
    names = append(names, "Dude", "Bro", "Pal", "Bub", "Man", "Brah", "Guy", "Buddy")

    client := &client{
        socket: socket,
        send:   make(chan []byte, messageBufferSize),
        room:   r,
        // Make random name here
        name:   fmt.Sprintf("%s%d", names[rand.Intn(len(names))], rand.Intn(100)),
    }

    r.join <- client
    defer func() { r.leave <- client }()

    go client.write()
    client.read()
}

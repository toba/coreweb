// Package socket enumerates HTTP headerSec-Websocket-* keys and handles
// websocket communications.
package socket

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/toba/coreweb/web"
)

type (
	// Request from browser as bytes along with WebSocket client the browser
	// communicated through.
	Request struct {
		Client  *Client
		Message []byte
	}

	// RequestHandler processes a socket request and returns a response that
	// should be sent to the client or nil if no response is expected.
	RequestHandler func(req *Request) []byte
)

const prefix = "Sec-Websocket-"

const (
	Accept   = prefix + "Accept"
	Key      = prefix + "Key"
	Protocol = prefix + "Protocol"
	Version  = prefix + "Version"
)

// Handle incoming websocket requests. Create a client object for each
// connection with a read and write event loop.
//
// Having pumps in goroutines allows "collection of memory referenced by the
// caller" according to
//
// https://github.com/gorilla/websocket/commit/ea4d1f681babbce9545c9c5f3d5194a789c89f5b
func Handle(c web.Config, responder RequestHandler) func(w http.ResponseWriter, r *http.Request) {
	broadcast = make(chan []byte)
	request = make(chan *Request)
	register = make(chan *Client)
	unregister = make(chan *Client)
	clients = make(map[*Client]bool)

	go listen(responder)

	// return standard HTTP handler that upgrades to socket connection
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
				debug.PrintStack()
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println(err)
			//http.Error(w, fmt.Sprintf("cannot upgrade: %v", err), http.StatusInternalServerError)
			return
		}
		client := &Client{conn: conn, Send: make(chan []byte, 256)}
		register <- client

		go client.writePump()
		go client.readPump()
	}
}

// listen is an event loop that continually checks event channels.
func listen(responder RequestHandler) {
	for {
		select {
		case c := <-register:
			clients[c] = true

		case c := <-unregister:
			if _, ok := clients[c]; ok {
				delete(clients, c)
				close(c.Send)
			}

		case req := <-request:
			res := responder(req)

			if res != nil {
				req.Client.Send <- res
			}

		case res := <-broadcast:
			for c := range clients {
				select {
				case c.Send <- res:
				default:
					close(c.Send)
					delete(clients, c)
				}
			}
		}
	}
}

// Broadcast puts a message onto the broadcast channel to be sent to all
// connected clients.
func Broadcast(res []byte) {
	if broadcast != nil && res != nil {
		broadcast <- res
	}
}

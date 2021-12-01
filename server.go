// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gorilla/websocket"
)

const (
	// Maximum message size allowed from peer.
	maxMessageSize = 512 * 512
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

var (
	addr     = flag.String("addr", ":8080", "http service address")
	newline  = []byte{'\n'}
	space    = []byte{' '}
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Connection to DB for updates
	db *bolt.DB
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

func main() {
	flag.Parse()
	db, err := initDB("einkauf.db")
	if err != nil {
		log.Fatal("BoltDB: ", err)
	}
	defer db.Close()
	hub := newHub(db)
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/api/v1/", func(w http.ResponseWriter, r *http.Request) {
		serveApi(w, r, db)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r, hub)
	})
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func initDB(pathToDB string) (*bolt.DB, error) {
	db, err := bolt.Open(pathToDB, 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("carts"))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("version"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method + " " + r.URL.String())
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	switch r.URL.Path {
	case "/":
		http.ServeFile(w, r, "public/index.html")
		return
	default:
		http.ServeFile(w, r, "public/"+r.URL.Path)
		return
	}
}

func serveApi(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	log.Println(r.Method + " " + r.URL.String())
	// Return the last 10 carts
	if r.Method == "GET" && r.URL.Path == "/api/v1/carts" {
		var data struct {
			Carts   map[string]interface{} `json:"carts"`
			Version string                 `json:"version"`
		}
		err := db.View(func(tx *bolt.Tx) error {
			version := tx.Bucket([]byte("version"))
			v1 := version.Get([]byte("cart"))
			data.Version = string(v1[:])
			c := tx.Bucket([]byte("carts")).Cursor()
			i := 1
			for k2, v2 := c.Last(); k2 != nil; k2, v2 = c.Prev() {
				var cart map[string]interface{}
				cartId := string(k2[:])
				if err := json.Unmarshal(v2, &cart); err != nil {
					return err
				}
				if i > 10 {
					break
				}
				carts, err := json.Marshal(map[string]interface{}{cartId: cart})
				if err != nil {
					return err
				}
				if err := json.Unmarshal(carts, &data.Carts); err != nil {
					return err
				}
				i++
			}
			result, err := json.Marshal(data)
			if err != nil {
				return err
			}
			_, err = w.Write(result)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			log.Printf("Error fetching data from API: %v", err)
			http.Error(w, "Error fetching data from API", http.StatusBadGateway)
			return
		}
		return
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request, hub *Hub) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		log.Println(r)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

func newHub(db *bolt.DB) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		db:         db,
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			var messageContent struct {
				Data struct {
					Cart    map[string]interface{} `json:"cart"`
					Version string                 `json:"version"`
				} `json:"data"`
				Tag string `json:"tag"`
			}
			if err := json.Unmarshal(message, &messageContent); err != nil {
				log.Printf("Error unmarshaling user data from websocket: %v", err)
			}
			// We received data from a websocket client that we weren't
			// expecting, ignore it
			if messageContent.Tag != "cartVersion" {
				return
			}
			cart, err := json.Marshal(messageContent.Data.Cart)
			if err != nil {
				log.Printf("Error unmarshaling user data from websocket: %v", err)
			}
			cartId := []byte(messageContent.Data.Version)
			err = h.db.Update(func(tx *bolt.Tx) error {
				carts := tx.Bucket([]byte("carts"))
				if err := carts.Put(cartId, cart); err != nil {
					return err
				}
				version := tx.Bucket([]byte("version"))
				if err := version.Put([]byte("cart"), cartId); err != nil {
					return err
				}
				return nil
			})
			log.Println("Sync changes to DB")
			if err != nil {
				log.Printf("Error sending data to the API: %v", err)
			}
		}
	}
}

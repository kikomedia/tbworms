package server

import (
	"log"
	"net/http"
	"tbworms/game"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

/*var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)*/

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	gameServer *GameServer
	conn       *websocket.Conn
	send       chan []byte
	player     *game.Player
	token      string
}

func (c *Client) readLoop() {
	defer func() {
		c.gameServer.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		log.Print("Received Msg: <" + string(message) + ">")
		if err != nil {
			log.Print("Error occured ", err.Error())

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//c.gameServer.onReceiveMessage(message, c)
		c.gameServer.command <- InputCommand{c, message}
		//c.gameServer.broadcast <- message
	}
}

func (c *Client) writeLoop() {
	pingTimer := time.NewTicker(pingPeriod)
	defer func() {
		pingTimer.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write([]byte(message))

			if err := w.Close(); err != nil {
				return
			}
		case <-pingTimer.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) run() {
	go c.writeLoop()
	go c.readLoop()
}

func (client Client) close() {
	client.conn.Close()
	close(client.send)
}

func newClient(g *GameServer, socket *websocket.Conn) *Client {
	return &Client{
		gameServer: g,
		conn:       socket,
		send:       make(chan []byte),
		player:     game.NewPlayer(),
		token:      uuid.NewV4().String(),
	}
}

func ServeWs(g *GameServer, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Upgrade Error", http.StatusInternalServerError)
		return
	}

	client := newClient(g, conn)
	g.register <- client
	g.clients = append(g.clients, client)
	client.run()
}

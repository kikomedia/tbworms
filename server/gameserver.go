package server

import (
	"encoding/json"
	"fmt"
	"time"

	message "tbworms/messages"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

type InputCommand struct {
	client  *Client
	command []byte
}

type GameServer struct {
	//clients    map[*Client]bool
	clients       []*Client
	broadcastChan chan []byte
	command       chan InputCommand
	register      chan *Client
	unregister    chan *Client
	//players       []game.Player

	controller *GameController
}

func NewGameServer() *GameServer {
	server := &GameServer{
		broadcastChan: make(chan []byte, 256),
		command:       make(chan InputCommand),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		clients:       make([]*Client, 0),
	}
	gameController := NewGameController(server)

	server.controller = gameController

	return server
}

func (g *GameServer) onInputCommand(input InputCommand) {
	fmt.Println("onInputCommand: " + string(input.command))

	message_type := gjson.GetBytes(input.command, "message_type").Int()
	if int(message_type) == message.MessageTypeUserInput {
		var msg message.UserKeyInput

		if json.Unmarshal(input.command, &msg) != nil {
			return
		}

		if msg.Token != input.client.token {
			fmt.Println("Invalid token")
		} else {
			g.controller.processPlayerInput(input.client.player, msg.Key)
		}
	}
}

func (g *GameServer) clientCount() int {
	return len(g.clients)
}

func (g *GameServer) onClientConnect(c *Client) {
	fmt.Println("onClientConnect..." + c.conn.RemoteAddr().String())
	g.broadcast(message.NewUserJoined(c.player), nil)
	g.sendToClient(message.NewOwnInfo(c.player.ID, c.token, c.player), c)
	g.sendToClient(message.NewGameAreaInfo(&g.controller.gameData.GameArea), c)
}

func (g *GameServer) sendToClient(message interface{}, client *Client) {
	data, _ := json.Marshal(message)
	//fmt.Println("ToClient: " + string(data))
	client.send <- data
}

func (g *GameServer) onClientDisconnect(c *Client) {
	fmt.Println("onClientDisconnect..." + c.conn.RemoteAddr().String())
	g.broadcast(message.NewUserLeft(c.player), c)

	c.close()
	// Find index of client
	i := -1
	for j, cl := range g.clients {
		if cl.token == c.token {
			i = j
			break
		}
	}
	copy(g.clients[i:], g.clients[i+1:])
	g.clients[len(g.clients)-1] = nil
	g.clients = g.clients[:len(g.clients)-1]
}

func (g *GameServer) broadcast(message interface{}, ignore *Client) {
	data, _ := json.Marshal(message)
	g.broadcastRaw(data, ignore)
}

func (g *GameServer) broadcastRaw(message []byte, ignore *Client) {
	//fmt.Println("Broadcast: " + string(message))
	for _, c := range g.clients {
		if c != ignore {
			c.send <- message
		}
	}
}

func (g *GameServer) Run() {
	step := time.NewTicker(GameTickerTime * time.Millisecond)
	defer func() {
		step.Stop()
	}()

	for {
		select {
		case <-step.C:
			g.controller.processStep()

		case client := <-g.register:
			g.onClientConnect(client)

		case client := <-g.unregister:
			g.onClientDisconnect(client)

		case input := <-g.command:
			g.onInputCommand(input)

		case message := <-g.broadcastChan:
			g.broadcastRaw(message, nil)
		}
	}
}

func (g *GameServer) Shutdown() {
	for _, client := range g.clients {
		client.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	}

	if g.clientCount() > 0 {
		time.Sleep(time.Millisecond * 2000)
	}
}

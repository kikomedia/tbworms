package server

import (
	"fmt"
	"math/rand"
	"strconv"
	"tbworms/game"
	message "tbworms/messages"
)

const PlayerStartPadding int = 2
const PlayerScoreStep int = 1
const PlayerScoreKill int = 10
const PlayerMaxRounds int = 10
const GameRoundWait int = 5000
const GameTickerTime = 100

type GameController struct {
	gameServer *GameServer
	gameData   *game.GameData
}

func NewGameController(gameServer *GameServer) *GameController {
	return &GameController{
		gameServer: gameServer,
		gameData:   game.NewGameData(),
	}
}

func (gc *GameController) processStep() {
	if gc.gameData.Waiting > 0 {
		gc.gameData.Waiting -= GameTickerTime
	} else {
		gc.processPlayerStep(false)
	}

	if gc.gameData.Waiting <= 0 {
		gc.gameData.Waiting = 0
		for _, cl := range gc.gameServer.clients {
			if cl.player.State == game.PlayerStateWaitingForNextRound ||
				cl.player.State == game.PlayerStateWon {
				cl.player.State = game.PlayerStatePlaying
			}
		}

		playing := 0
		died := 0

		for _, cl := range gc.gameServer.clients {
			switch cl.player.State {
			case game.PlayerStatePlaying:
				playing += 1
			case game.PlayerStateDied:
				died += 1
			}
		}

		if (playing == 0) || (playing == 1 && died > 0) {
			if playing == 1 && died > 0 {
				for _, cl := range gc.gameServer.clients {
					if cl.player.State == game.PlayerStatePlaying {
						cl.player.State = game.PlayerStateWon
					}
				}
			}

			gc.startRound()

		}

	}

	// - Send data to clients

	players := []game.Player{}

	for _, cl := range gc.gameServer.clients {
		players = append(players, *cl.player)
	}

	gc.gameServer.broadcast(message.NewServerUpdate(players, gc.gameData.GameState, gc.gameData.Waiting, gc.gameData.Round), nil)
}

func (gc *GameController) movePlayer(player *game.Player, x int, y int) {
	px := player.HeadPos.X
	py := player.HeadPos.Y

	px = px + x
	py = py + y

	if px > gc.gameData.GameArea.Width-1 ||
		py > gc.gameData.GameArea.Height-1 ||
		px < 0 ||
		py < 0 {
		fmt.Println("Player ", player.Username, " died")
		player.State = game.PlayerStateDied

	} else {
		player.HeadPos.X = px
		player.HeadPos.Y = py

		gc.checkCollision(player, px, py)

		player.Positions = append(player.Positions, game.Position{X: px, Y: py})
	}

	//fmt.Println("Player:"+player.Username+" HeadPos: ", player.HeadPos, " Pos :", player.Positions)
}

func (gc *GameController) checkCollision(player *game.Player, x int, y int) {
	px := player.HeadPos.X
	py := player.HeadPos.Y

	for _, cl := range gc.gameServer.clients {
		for _, p := range cl.player.Positions {
			if p.X == px && p.Y == py {
				cl.player.Score += PlayerScoreKill
				player.State = game.PlayerStateDied
				return
			}
		}
	}
}

func (gc *GameController) processPlayerStep(init_step bool) {
	for _, cl := range gc.gameServer.clients {

		if (cl.player.State == game.PlayerStatePlaying) ||
			(init_step && cl.player.State == game.PlayerStateWaitingForNextRound) ||
			(init_step && cl.player.State == game.PlayerStateWon) {
			switch cl.player.Direction {
			case game.PlayerDirectionLeft:
				gc.movePlayer(cl.player, -1, 0)
			case game.PlayerDirectionRight:
				gc.movePlayer(cl.player, 1, 0)
			case game.PlayerDirectionUp:
				gc.movePlayer(cl.player, 0, -1)
			case game.PlayerDirectionDown:
				gc.movePlayer(cl.player, 0, 1)
			default:
				fmt.Println("No Movement " + strconv.Itoa(cl.player.Direction))
			}
		}
	}

	for _, cl := range gc.gameServer.clients {
		switch cl.player.State {
		case game.PlayerStatePlaying:
			cl.player.Score += PlayerScoreStep
		case game.PlayerStateDied:
			cl.player.Positions = []game.Position{}
		default:
			//
		}
	}
}

func (gc *GameController) processPlayerInput(player *game.Player, key int) {
	fmt.Println("User ", player.Username, " pressed :  ", key)

	if player.State == game.PlayerStatePlaying {

		switch key {
		case message.UserInputKeyLeft:
			player.Direction = game.PlayerDirectionLeft
		case message.UserInputKeyRight:
			player.Direction = game.PlayerDirectionRight
		case message.UserInputKeyUp:
			player.Direction = game.PlayerDirectionUp
		case message.UserInputKeyDown:
			player.Direction = game.PlayerDirectionDown
		default:
			fmt.Println("Invalid input")
		}
	}
}

func (gc *GameController) resetPlayerDataForNewRound(player *game.Player, reset_score bool) {
	player.HeadPos = gc.GenerateNewHeadPos(
		PlayerStartPadding,
		gc.gameData.GameArea.Width-PlayerStartPadding,
		PlayerStartPadding,
		gc.gameData.GameArea.Height-PlayerStartPadding)

	player.Positions = []game.Position{}

	player.Direction = rand.Intn(3) + 1
	if reset_score {
		player.Score = 0
	}
}

func (gc *GameController) startRound() {
	reset_score := false
	gc.gameData.Round += 1
	gc.gameData.Waiting = GameRoundWait

	if gc.gameData.Round > PlayerMaxRounds {
		gc.gameData.Round = 1
		reset_score = true
	}

	for _, cl := range gc.gameServer.clients {
		if cl.player.State == game.PlayerStateWaitingForNextRound ||
			cl.player.State == game.PlayerStateConnected ||
			cl.player.State == game.PlayerStateDied ||
			cl.player.State == game.PlayerStatePlaying {
			cl.player.State = game.PlayerStateWaitingForNextRound
			gc.resetPlayerDataForNewRound(cl.player, reset_score)
		}

		if cl.player.State == game.PlayerStateWon {
			gc.resetPlayerDataForNewRound(cl.player, reset_score)
		}
	}
	gc.processPlayerStep(true)
	gc.processPlayerStep(true)
}
func (p *GameController) GenerateNewHeadPos(x1 int, x2 int, y1 int, y2 int) game.Position {
	return game.Position{
		X: rand.Intn(x2-x1) + x1,
		Y: rand.Intn(y2-y1) + y1,
	}
}

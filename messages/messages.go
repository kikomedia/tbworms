package message

import "tbworms/game"

const (
	MessageTypeUserJoined   int = 1000
	MessageTypeUserLeft     int = 1001
	MessageTypeOwnInfo      int = 2000
	MessageTypeServerUpdate int = 3000
	MessageTypeGameAreaInfo int = 3001
	MessageTypeNewRound     int = 3002
	MessageTypeUserInput    int = 5000
)

const (
	UserInputKeyLeft   int = 37
	UserInputKeyRight  int = 39
	UserInputKeyUp     int = 38
	UserInputKeyDown   int = 40
	UserInputKeyEnter  int = 13
	UserInputKeyEscape int = 27
	UserInputKeySpace  int = 32
)

type MessageType struct {
	MessageType int `json:"message_type"`
}

type UserJoined struct {
	MessageType int         `json:"message_type"`
	Player      game.Player `json:"player"`
}

type UserKeyInput struct {
	MessageType int    `json:"message_type"`
	Token       string `json:"token"`
	Key         int    `json:"key"`
}

type GameAreaInfo struct {
	MessageType int           `json:"message_type"`
	GameArea    game.GameArea `json:"game_area"`
}

type OwnInfo struct {
	MessageType int         `json:"message_type"`
	OwnID       string      `json:"own_id"`
	Token       string      `json:"token"`
	Player      game.Player `json:"player"`
}

type UserLeft struct {
	MessageType int         `json:"message_type"`
	Player      game.Player `json:"player"`
}

type ServerUpdate struct {
	MessageType int           `json:"message_type"`
	State       int           `json:"state"`
	Waiting     int           `json:"waiting"`
	Round       int           `json:"round"`
	Info        string        `json:"info"`
	Players     []game.Player `json:"players"`
}

func NewUserJoined(player *game.Player) *UserJoined {
	return &UserJoined{
		MessageType: MessageTypeUserJoined,
		Player:      *player,
	}
}

func NewUserLeft(player *game.Player) *UserJoined {
	return &UserJoined{
		MessageType: MessageTypeUserLeft,
		Player:      *player,
	}
}

func NewOwnInfo(own_id string, token string, player *game.Player) *OwnInfo {
	return &OwnInfo{
		MessageType: MessageTypeOwnInfo,
		OwnID:       own_id,
		Token:       token,
		Player:      *player,
	}
}

func NewServerUpdate(players []game.Player, state int, waiting int, round int) *ServerUpdate {
	return &ServerUpdate{
		MessageType: MessageTypeServerUpdate,
		State:       state,
		Waiting:     waiting,
		Players:     players,
		Round:       round,
	}
}

func NewGameAreaInfo(gameArea *game.GameArea) *GameAreaInfo {
	return &GameAreaInfo{
		MessageType: MessageTypeGameAreaInfo,
		GameArea:    *gameArea,
	}
}

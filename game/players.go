package game

import (
	"math/rand"
	"tbworms/utils"

	uuid "github.com/satori/go.uuid"
)

const (
	PlayerDirectionRight = iota
	PlayerDirectionUp
	PlayerDirectionLeft
	PlayerDirectionDown
)

const (
	PlayerStateDisconnected = iota
	PlayerStateConnected
	PlayerStateWaitingForNextRound
	PlayerStatePlaying
	PlayerStateDied
	PlayerStatePaused
	PlayerStateWon
)

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PositionList struct {
	Points []Position `json:"points"`
}

type Player struct {
	ID        string     `json:"id"`
	Color     string     `json:"color"`
	Username  string     `json:"username"`
	State     int        `json:"state"`
	Direction int        `json:"direction"`
	HeadPos   Position   `json:"headpos"`
	Score     int        `json:"score"`
	Positions []Position `json:"positions"`
}

func NewPlayer() *Player {
	return &Player{
		State:     PlayerStateConnected,
		ID:        uuid.NewV4().String(),
		Color:     utils.GenerateColor(),
		Username:  utils.GenerateUsername(),
		Direction: rand.Intn(3),
		HeadPos:   Position{X: 0, Y: 0},
		Score:     0,
		Positions: []Position{},
	}
}

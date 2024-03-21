package game

const (
	GameStateWaiting int = iota
	GameStatePlaying
	GameStateFinished
)

type GameData struct {
	GameArea  GameArea
	GameState int
	Round     int
	Waiting   int
}

func NewGameData() *GameData {
	return &GameData{
		GameArea:  NewGameArea(40, 40),
		GameState: GameStateWaiting,
	}
}

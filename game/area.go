package game

type GameArea struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func NewGameArea(width int, height int) GameArea {
	return GameArea{Width: width, Height: height}
}

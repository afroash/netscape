package entities

import "github.com/hajimehoshi/ebiten"

type Player struct {
	PlayerImage      *ebiten.Image
	PlayerX, PlayerY float64
}

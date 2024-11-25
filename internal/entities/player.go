package entities

import "github.com/hajimehoshi/ebiten/v2"

type Player struct {
	PlayerImage      *ebiten.Image
	PlayerX, PlayerY float64
}

package main

import (
	"log"
	"os"

	"github.com/afroash/netscape/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("NetScape!")
	//ebiten.SetWindowResizable(true)
	game := game.NewGame()
	if err := ebiten.RunGame(game); err != nil {
		if err == ebiten.Termination {
			//Clean Exit
			os.Exit(0)
		}

		log.Fatal(err)
	}

}

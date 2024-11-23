package main

import (
	"image"
	"image/color"
	"log"

	"github.com/afroash/netscape/internal/camera"
	"github.com/afroash/netscape/internal/drawstuff"
	"github.com/afroash/netscape/internal/entities"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Game struct {
	Player       *entities.Player
	TileMapJson  *drawstuff.TileMapJson
	tileMapImage *ebiten.Image
	cam          *camera.Camera
}

func (g *Game) Update(*ebiten.Image) error {
	// Player moves Left
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.Player.PlayerX -= 2
	}
	// Player moves Right
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.Player.PlayerX += 2
	}
	// Player moves Up
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.Player.PlayerY -= 2
	}
	// Player moves Down
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.Player.PlayerY += 2
	}

	// Camera follows the player
	g.cam.FollowPlayer(g.Player.PlayerX+8, g.Player.PlayerY+8, 320, 240)
	// Camera is constrained to the map
	g.cam.Constrain(
		320.0,
		240.0,
		float64(g.TileMapJson.Layers[0].Width)*16.0,
		float64(g.TileMapJson.Layers[0].Height)*16.0,
	)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the game screen here
	screen.Fill(color.Gray{0x80})
	op := &ebiten.DrawImageOptions{}
	// Draw the tilemap
	for _, layer := range g.TileMapJson.Layers {
		for i, id := range layer.Data {
			x := float64((i % layer.Width) * 16)
			y := float64((i / layer.Width) * 16)

			srcX := ((id - 1) % 16) * 16
			srcY := ((id - 1) / 16) * 16

			op.GeoM.Translate(x, y)
			op.GeoM.Translate(g.cam.CameraX, g.cam.CameraY)

			screen.DrawImage(g.tileMapImage.SubImage(image.Rect(
				srcX, srcY, srcX+16, srcY+16,
			)).(*ebiten.Image), op)

			op.GeoM.Reset()
		}
	}

	op.GeoM.Translate(g.Player.PlayerX, g.Player.PlayerY)
	op.GeoM.Translate(g.cam.CameraX, g.cam.CameraY)
	screen.DrawImage(g.Player.PlayerImage.SubImage(image.Rect(
		0, 104, 16, 128,
	)).(*ebiten.Image), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	//Return the game screen size here
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("NetScape!")
	ebiten.SetWindowResizable(true)
	playerImg, _, err := ebitenutil.NewImageFromFile("assests/images/player.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	tileMapImage, _, err := ebitenutil.NewImageFromFile("assests/images/PixelOffice.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	tileMapJson, err := drawstuff.NewTileMapJson("assests/Maps/Floors/floor1.json")
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{
		Player: &entities.Player{
			PlayerImage: playerImg,
			PlayerX:     100,
			PlayerY:     100,
		},
		TileMapJson:  tileMapJson,
		tileMapImage: tileMapImage,
		cam:          camera.NewCamera(0.0, 0.0),
	}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}

	//Debugging area print while game is running.
	// fmt.Println("Player X: ", game.Player.PlayerX)

}

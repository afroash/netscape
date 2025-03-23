package game

import (
	"bytes"
	"image"
	"image/color"
	"log"

	"github.com/afroash/netscape/internal/camera"
	debug "github.com/afroash/netscape/internal/debugy"
	"github.com/afroash/netscape/internal/drawstuff"
	"github.com/afroash/netscape/internal/interaction"

	"github.com/afroash/netscape/internal/entities"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type GameState int

const (
	MainMenu GameState = iota
	Playing
)

type Game struct {
	Player            *entities.Player
	TileMapJson       *drawstuff.TileMapJson
	TileMapImage      *ebiten.Image
	Cam               *camera.Camera
	ShouldExit        bool
	GameState         GameState
	Menu              *drawstuff.Menu
	DrawStuff         *drawstuff.DrawStuff
	InteractionPoints []*interaction.InteractionPoint
	DialogeBox        *interaction.DialogeBox
	Debugy            *debug.DebugInfo
}

func NewGame() *Game {
	// Initialize menu
	menu := &drawstuff.Menu{
		Items: []drawstuff.MenuItem{
			{Text: "New Game", Selected: true},
			{Text: "Exit", Selected: false},
		},
		Selected: 0,
	}

	// Initialize DrawStuff with font
	drawStuff, err := drawstuff.NewDrawStuff(menu)
	if err != nil {
		log.Fatal("Error initializing DrawStuff:", err)
	}

	//initailize font for dialoge box
	diaglogFontSource, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal("Error initializing DialogeBox font:", err)
	}

	diaglogFont := &text.GoTextFace{
		Source: diaglogFontSource,
		Size:   8,
	}

	game := &Game{
		GameState: MainMenu,
		Menu:      menu,
		DrawStuff: drawStuff,
		InteractionPoints: []*interaction.InteractionPoint{
			{
				X:     170,
				Y:     140,
				Range: 30,
				Messages: []string{
					"Welcome to the office! This is your first day.",
					"Your desk is located in the corner.",
					"Press 'E' to interact with objects when prompted.",
				},
				IsActive:      true,
				HasInteracted: false,
				CurrentMsg:    1,
			},
			{
				X:     150,
				Y:     400,
				Range: 30,
				Messages: []string{
					"This is the break room.",
					"Take breaks regularly to maintain productivity!",
				},
				IsActive:      true,
				HasInteracted: false,
				CurrentMsg:    2,
			},
		},
		DialogeBox: interaction.NewDialogeBox(diaglogFont),
		Debugy:     debug.NewDebugInfo(diaglogFont),
	}

	return game
}

func (g *Game) initializeGameResources() error {
	// Only load game resources when starting a new game
	playerImg, _, err := ebitenutil.NewImageFromFile("assests/images/player.png")
	if err != nil {
		return err
	}

	tileMapImage, _, err := ebitenutil.NewImageFromFile("assests/images/PixelOffice.png")
	if err != nil {
		return err
	}

	tileMapJson, err := drawstuff.NewTileMapJson("assests/Maps/Floors/floor1.json")
	if err != nil {
		return err
	}

	g.Player = &entities.Player{
		PlayerImage: playerImg,
		PlayerX:     100,
		PlayerY:     100,
	}
	g.TileMapJson = tileMapJson
	g.TileMapImage = tileMapImage
	g.Cam = camera.NewCamera(0.0, 0.0)

	return nil
}

func (g *Game) Update() error {
	if g.ShouldExit {
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		g.Debugy.Enabled = !g.Debugy.Enabled
	}

	switch g.GameState {
	case MainMenu:
		// Handle menu navigation
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			g.Menu.Selected = (g.Menu.Selected - 1 + len(g.Menu.Items)) % len(g.Menu.Items)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			g.Menu.Selected = (g.Menu.Selected + 1) % len(g.Menu.Items)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			switch g.Menu.Selected {
			case 0: // New Game
				if err := g.initializeGameResources(); err != nil {
					return err
				}
				g.GameState = Playing
			case 1: // Exit
				g.ShouldExit = true
			}
		}

	case Playing:
		// Existing game update logic
		// Only process interactions when in playing state
		if !g.DialogeBox.IsVisible {
			// Handle player movement only when dialog is not showing
			if ebiten.IsKeyPressed(ebiten.KeyLeft) {
				g.Player.PlayerX -= 2
			}
			if ebiten.IsKeyPressed(ebiten.KeyRight) {
				g.Player.PlayerX += 2
			}
			if ebiten.IsKeyPressed(ebiten.KeyUp) {
				g.Player.PlayerY -= 2
			}
			if ebiten.IsKeyPressed(ebiten.KeyDown) {
				g.Player.PlayerY += 2
			}

			if ebiten.IsKeyPressed(ebiten.KeyZ) {
				return ebiten.Termination
			}

			// Check for interactions
			for _, point := range g.InteractionPoints {
				if point.IsPlayerInRange(g.Player.PlayerX, g.Player.PlayerY) {
					if !point.IsActive {
						point.IsActive = true
						// Show 'Press E to interact' indicator
					}

					if inpututil.IsKeyJustPressed(ebiten.KeyE) && !point.HasInteracted {
						// Only mark the current point as interacted
						point.HasInteracted = true
						point.CurrentMsg = 0
						g.DialogeBox.IsVisible = true
						g.DialogeBox.CurrentPoint = point
						break // Exit the loop after interacting with one point
					}
				} else {
					point.IsActive = false
				}
			}
		} else {
			// Handle dialog navigation
			if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
				currentPoint := g.DialogeBox.CurrentPoint
				if currentPoint.CurrentMsg < len(currentPoint.Messages)-1 {
					currentPoint.CurrentMsg++
				} else {
					g.DialogeBox.IsVisible = false
					g.DialogeBox.CurrentPoint = nil
				}
			}
		}

		// Update camera
		g.Cam.FollowPlayer(g.Player.PlayerX+8, g.Player.PlayerY+8, 320, 240)
		g.Cam.Constrain(
			320.0,
			240.0,
			float64(g.TileMapJson.Layers[0].Width)*16.0,
			float64(g.TileMapJson.Layers[0].Height)*16.0,
		)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xf0, 0xf0, 0xf0, 0xff})

	switch g.GameState {
	case MainMenu:
		drawstuff.DrawMenu(screen, g.DrawStuff)

	case Playing:
		// Existing game drawing logic
		for _, layer := range g.TileMapJson.Layers {
			op := &ebiten.DrawImageOptions{}
			for i, id := range layer.Data {
				x := float64((i % layer.Width) * 16)
				y := float64((i / layer.Width) * 16)

				srcX := ((id - 1) % 16) * 16
				srcY := ((id - 1) / 16) * 16

				op.GeoM.Reset()
				op.GeoM.Translate(x, y)
				op.GeoM.Translate(g.Cam.CameraX, g.Cam.CameraY)

				screen.DrawImage(g.TileMapImage.SubImage(image.Rect(
					srcX, srcY, srcX+16, srcY+16,
				)).(*ebiten.Image), op)
			}
		}

		// Draw player
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(g.Player.PlayerX, g.Player.PlayerY)
		op.GeoM.Translate(g.Cam.CameraX, g.Cam.CameraY)
		screen.DrawImage(g.Player.PlayerImage.SubImage(image.Rect(
			0, 104, 16, 128,
		)).(*ebiten.Image), op)

		// Draw interaction points
		for _, point := range g.InteractionPoints {
			if point.IsActive && !point.HasInteracted {
				// Draw 'Press E to interact' above the interaction point
				op := &text.DrawOptions{}
				worldX := point.X + g.Cam.CameraX
				worldY := point.Y + g.Cam.CameraY - 20
				op.GeoM.Translate(worldX, worldY)
				op.ColorScale.ScaleWithColor(color.White)
				text.Draw(screen, "Press E", g.DialogeBox.FontFace, op)
			}
		}
		g.DialogeBox.Draw(screen)

		// Draw debug info
		if g.Debugy.Enabled {
			debugParams := map[string]interface{}{
				"x":    g.Player.PlayerX,
				"y":    g.Player.PlayerY,
				"camX": g.Cam.CameraX,
				"camY": g.Cam.CameraY,
				"fps":  ebiten.ActualFPS(),
			}
			g.Debugy.Draw(screen, debugParams)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

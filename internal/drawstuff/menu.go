package drawstuff

import (
	"bytes"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth    = 320
	screenHeight   = 240
	menuFontSize   = 24
	diffFontSize   = 20
	normalFontSize = 12
)

type Menu struct {
	Items    []MenuItem
	Selected int
}

type MenuItem struct {
	Text     string
	Selected bool
}

type DrawStuff struct {
	fontFace     *text.GoTextFace
	titleFace    *text.GoTextFace
	instructFace *text.GoTextFace
	game         *Menu
}

// NewDrawStuff creates a new DrawStuff instance with initialized fonts
func NewDrawStuff(menu *Menu) (*DrawStuff, error) {

	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	// Create different font faces for different sizes
	menuFont := &text.GoTextFace{
		Source: fontSource,
		Size:   menuFontSize,
	}

	titleFont := &text.GoTextFace{
		Source: fontSource,
		Size:   menuFontSize + 8,
	}

	instructFont := &text.GoTextFace{
		Source: fontSource,
		Size:   normalFontSize,
	}

	return &DrawStuff{
		fontFace:     menuFont,
		titleFace:    titleFont,
		instructFace: instructFont,
		game:         menu,
	}, nil
}

func DrawMenu(screen *ebiten.Image, d *DrawStuff) {
	if d == nil || d.fontFace == nil {
		log.Println("Warning: DrawStuff or font not initialized")
		return
	}

	// Center the menu on screen
	startX := screenWidth / 2
	startY := screenHeight / 2

	// Draw title
	titleOp := &text.DrawOptions{}
	titleOp.GeoM.Translate(float64(startX), float64(startY-85))
	titleOp.ColorScale.ScaleWithColor(color.Black)
	titleOp.PrimaryAlign = text.AlignCenter
	text.Draw(screen, "Netscape by Ash", d.titleFace, titleOp)

	lineSpacing := 40 // Increased spacing between options
	options := []string{"New Game", "Exit"}

	for i, option := range options {
		yPos := startY - 10 + i*lineSpacing

		// Calculate text metrics for centering
		textWidth := len(option) * diffFontSize / 2 // Approximate width
		rectWidth := float32(textWidth + 40)        // Add padding
		rectHeight := float32(40)                   // Fixed height for selection rectangle

		// Draw selection highlight if this option is selected
		if i == d.game.Selected {
			vector.DrawFilledRect(
				screen,
				float32(startX)-rectWidth/2,
				float32(yPos)-rectHeight/2,
				rectWidth,
				rectHeight,
				color.RGBA{0, 0, 255, 100},
				false,
			)

			vector.StrokeRect(
				screen,
				float32(startX)-rectWidth/2,
				float32(yPos)-rectHeight/2,
				rectWidth,
				rectHeight,
				2,
				color.RGBA{0, 0, 255, 255},
				false,
			)
		}

		// Draw the menu option text
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(startX), float64(yPos))
		op.ColorScale.ScaleWithColor(color.Black)
		op.PrimaryAlign = text.AlignCenter
		op.SecondaryAlign = text.AlignCenter

		text.Draw(screen, option, d.fontFace, op)
	}

	// Draw instructions at the bottom
	instructOp := &text.DrawOptions{}
	instructOp.GeoM.Translate(float64(startX), float64(startY+lineSpacing*2))
	instructOp.ColorScale.ScaleWithColor(color.RGBA{100, 100, 100, 255})
	instructOp.PrimaryAlign = text.AlignCenter
	text.Draw(screen, "Use ↑↓ to select, ENTER to confirm, 'z' to exit game", d.instructFace, instructOp)
}

package interaction

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type InteractionPoint struct {
	X, Y          float64
	Range         float64
	Messages      []string
	CurrentMsg    int
	IsActive      bool
	HasInteracted bool
}

// Dialogebox handles the rendring of the dialoge box
type DialogeBox struct {
	Width, Height int
	Padding       int
	IsVisible     bool
	FontFace      *text.GoTextFace
	CurrentPoint  *InteractionPoint
}

// NewDialogeBox creates a new DialogeBox instance
func NewDialogeBox(fontFace *text.GoTextFace) *DialogeBox {
	return &DialogeBox{
		Width:     280,
		Height:    80,
		Padding:   10,
		IsVisible: false,
		FontFace:  fontFace,
	}
}

// IsPlayerinRange checks if the player is in range of the interaction point
func (ip *InteractionPoint) IsPlayerInRange(playerX, playerY float64) bool {
	dx := ip.X - playerX
	dy := ip.Y - playerY
	distance := math.Sqrt(dx*dx + dy*dy)
	return distance <= ip.Range*ip.Range
}

// Draw renders the dialoge box
func (db *DialogeBox) Draw(screen *ebiten.Image) {
	if !db.IsVisible || db.CurrentPoint == nil {
		return
	}

	screenWidth, screenHeight := screen.Bounds().Dx(), screen.Bounds().Dy()

	// Position the dialog box at the bottom center of the screen
	x := (screenWidth - db.Width) / 2
	y := screenHeight - db.Height - 20

	// Draw semi-transparent background
	vector.DrawFilledRect(
		screen,
		float32(x),
		float32(y),
		float32(db.Width),
		float32(db.Height),
		color.RGBA{0, 0, 0, 200},
		true,
	)

	// Draw border
	vector.StrokeRect(
		screen,
		float32(x),
		float32(y),
		float32(db.Width),
		float32(db.Height),
		2,
		color.RGBA{255, 255, 255, 255},
		true,
	)

	// Draw text
	if db.CurrentPoint.CurrentMsg < len(db.CurrentPoint.Messages) {
		message := db.CurrentPoint.Messages[db.CurrentPoint.CurrentMsg]

		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(x+db.Padding), float64(y+db.Padding))
		op.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, message, db.FontFace, op)

		// Draw "Press Space to continue" if there are more messages
		if db.CurrentPoint.CurrentMsg < len(db.CurrentPoint.Messages)-1 {
			continueOp := &text.DrawOptions{}
			continueOp.GeoM.Translate(
				float64(x+db.Width-db.Padding-100),
				float64(y+db.Height-db.Padding-15),
			)
			continueOp.ColorScale.ScaleWithColor(color.RGBA{200, 200, 200, 255})
			text.Draw(screen, "Press Space →", db.FontFace, continueOp)
		}
	}
}

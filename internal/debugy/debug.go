package debug

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type DebugInfo struct {
	Enabled  bool
	FontFace *text.GoTextFace
}

func NewDebugInfo(fontFace *text.GoTextFace) *DebugInfo {
	return &DebugInfo{
		Enabled:  true,
		FontFace: fontFace,
	}
}

func (d *DebugInfo) Draw(screen *ebiten.Image, params map[string]interface{}) {
	if !d.Enabled {
		return
	}

	// Semi-transparent black background for debug info
	debugBgColor := color.RGBA{0, 0, 0, 180}
	debugTextColor := color.RGBA{0, 255, 0, 255} // Green text for debug info

	// Draw debug background
	width := float32(screen.Bounds().Dx())
	height := float32(25) // Height for debug bar
	vector.DrawFilledRect(screen, 0, 0, width, height, debugBgColor, true)

	// Draw debug text
	textY := float64(15)    // Vertical position for text
	spacing := float64(120) // Horizontal spacing between info segments

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(debugTextColor)

	// Draw each debug parameter
	for i, pair := range formatDebugParams(params) {
		op.GeoM.Reset()
		op.GeoM.Translate(float64(10+(float64(i)*spacing)), textY)
		text.Draw(screen, pair, d.FontFace, op)
	}
}

func formatDebugParams(params map[string]interface{}) []string {
	formatted := make([]string, 0)

	// Format position if available
	if x, ok := params["x"]; ok {
		if y, ok := params["y"]; ok {
			formatted = append(formatted, fmt.Sprintf("Pos: %.1f, %.1f", x, y))
		}
	}

	// Format camera if available
	if camX, ok := params["camX"]; ok {
		if camY, ok := params["camY"]; ok {
			formatted = append(formatted, fmt.Sprintf("Cam: %.1f, %.1f", camX, camY))
		}
	}

	// Add any other debug info
	if fps, ok := params["fps"]; ok {
		formatted = append(formatted, fmt.Sprintf("FPS: %.1f", fps))
	}

	return formatted
}

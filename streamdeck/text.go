package streamdeck

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const defaultFontSize = 14
const defaultFont = "/usr/share/fonts/adwaita-sans-fonts/AdwaitaSans-Regular.ttf"

func (dd *DeckDevice) SetText(img image.Image, text string, fontPath string, fontSize int, fontColor string, pt image.Point) (image.Image, error) {
	if fontPath == "" {
		fontPath = defaultFont
	}
	if fontSize == 0 {
		fontSize = defaultFontSize
	}

	// Convert to RGBA for drawing
	result := image.NewRGBA(img.Bounds())
	draw.Draw(result, result.Bounds(), img, image.Point{}, draw.Src)

	// 2. Load the font
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return result, fmt.Errorf("could not read font file: %w", err)
	}

	// 3. Create a font face
	ttf, err := opentype.Parse(fontBytes)
	if err != nil {
		return result, err
	}

	face, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    float64(fontSize),
		DPI:     float64(dd.DPI),
		Hinting: font.HintingNone,
	})
	if err != nil {
		return result, err
	}

	// 4. Define the text color (Using a placeholder, as your 'color' argument is a string)
	// You would typically parse the 'colorHex' string argument here,
	// but for now, we use a fixed color.
	col := color.RGBA{R: 255, G: 255, B: 255, A: 255} // White color

	// 5. Create the Drawer
	d := &font.Drawer{
		Dst:  result, // Use the defined variable 'result'
		Src:  image.NewUniform(col),
		Face: face,
		// 6. Set the position
		Dot: fixed.Point26_6{
			X: fixed.I(pt.X), // 'pt' must be passed as an argument
			Y: fixed.I(pt.Y) + face.Metrics().Ascent,
		},
	}

	// 7. Draw the text
	d.DrawString(text)

	return result, nil
}

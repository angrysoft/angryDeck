package streamdeck

import (
	"fmt"
	"image"
	"image/color"
	imgDraw "image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"golang.org/x/image/draw"
)

func (dd *DeckDevice) SetImageFromFile(keyIndex uint8, path string) error {
	img, err := loadImageFromFile(path)
	if err != nil {
		return err
	}
	return dd.SetImage(keyIndex, img)
}

func (dd *DeckDevice) SetImageWithTextFromFile(keyIndex uint8, path string, text string) error {
	img, err := loadImageFromFile(path)
	if err != nil {
		return err
	}
	img, err = dd.SetText(img, text, "", 0, "white", image.Point{X: 0, Y: 0})
	if err != nil {
		return err
	}
	return dd.SetImage(keyIndex, img)
}

func (dd *DeckDevice) SetImage(keyIndex uint8, img image.Image) error {
	imageData, err := dd.prepareImage(img)
	if err != nil {
		return err
	}
	data := make([]byte, dd.imagePageSize)
	translatedIndex := dd.translateKeyIndex(keyIndex, dd.Columns)

	var page int
	var lastPage bool
	for !lastPage {
		var payload []byte
		payload, lastPage = imageData.page(page)
		header := dd.imagePageHeader(page, translatedIndex, len(payload), lastPage)

		copy(data, header)
		copy(data[len(header):], payload)

		err := dd.Write(data)
		if err != nil {
			return fmt.Errorf("cannot write image page %d (%d bytes): %v",
				page, len(data), err)
		}

		page++
	}

	return nil
}

func createNewRGBAImage(width, height int, fillColor color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	imgDraw.Draw(img, img.Bounds(), image.NewUniform(fillColor), image.Point{}, imgDraw.Src)
	return img
}

func loadImageFromFile(path string) (image.Image, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open image file: %v", err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("cannot decode image file: %v", err)
	}

	return img, nil
}

func resizeImage(src image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), imgDraw.Over, nil)
	return dst
}

func (dd *DeckDevice) prepareImage(img image.Image) (*ImageData, error) {
	// Auto-resize if dimensions don't match
	if img.Bounds().Dy() != int(dd.Pixels) || img.Bounds().Dx() != int(dd.Pixels) {
		img = resizeImage(img, int(dd.Pixels), int(dd.Pixels))
	}

	imageBytes, err := dd.toImageFormat(dd.flipImage(img))
	if err != nil {
		return nil, fmt.Errorf("cannot convert image data: %v", err)
	}

	pageSize := dd.imagePageSize - dd.imagePageHeaderSize

	pageCount := len(imageBytes) / pageSize
	if len(imageBytes)%pageSize != 0 {
		pageCount++
	}

	return &ImageData{
		image:     imageBytes,
		pageSize:  pageSize,
		pageCount: pageCount,
	}, nil
}

// Clears the Stream Deck, setting a black image on all buttons.
func (dd *DeckDevice) Clear() error {
	img := image.NewRGBA(image.Rect(0, 0, int(dd.Pixels), int(dd.Pixels)))
	imgDraw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{0, 0, 0, 255}), image.Point{}, imgDraw.Src)
	for i := uint8(0); i <= dd.Columns*dd.Rows; i++ {
		err := dd.SetImage(i, img)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

type ImageData struct {
	image     []byte
	pageSize  int
	pageCount int
}

// page returns the page with the given index and an indication if this is the
// last page.
func (d ImageData) page(pageIndex int) ([]byte, bool) {
	offset := pageIndex * d.pageSize
	length := d.pageSize

	if offset+length > len(d.image) {
		length = len(d.image) - offset
	}

	if length <= 0 {
		return []byte{}, true
	}

	return d.image[offset : offset+length], pageIndex == d.pageCount-1
}

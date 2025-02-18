package icat

import (
	"bytes"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/BourgeoisBear/rasterm"
	"github.com/hilli/icat/util"
	"github.com/nfnt/resize"

	ascii "github.com/qeesung/image2ascii/convert"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/riff"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/vp8"
	_ "golang.org/x/image/vp8l"
	_ "golang.org/x/image/webp"
)

func PrintImageFile(imageFileName string, forceASCII bool) error {
	imageFile, imageSize, err := util.FileAndStat(imageFileName)
	if err != nil {
		return err
	}
	defer imageFile.Close()

	imageData, err := os.ReadFile(imageFileName)
	if err != nil {
		return err
	}

	img, err := DecodeImage(imageData)
	if err != nil {
		return err
	}
	return PrintImage(img, imageFileName, imageSize, forceASCII)
}

func PrintImageURL(imageURL string, forceASCII bool) error {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(imageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	img, err := DecodeImage(imageData)
	if err != nil {
		return err
	}
	return PrintImage(img, imageURL, resp.ContentLength, forceASCII)
}

func PrintImage(img *image.Image, filename string, imageSize int64, forceASCII bool) error {
	var img2 image.Image = *img

	if forceASCII {
		fmt.Print("\n", ConvertToASCII(img2))
		return nil
	}

	sixelCapable, _ := rasterm.IsSixelCapable()

	_, _, pw, ph := TermSize() // Get terminal height and width in pixels
	size, resizeOption := resizeConstraints(img2.Bounds(), int(pw), int(ph))

	switch resizeOption {
	case 'x':
		img2 = resize.Resize(uint(size), 0, img2, resize.NearestNeighbor)
	case 'y':
		img2 = resize.Resize(0, uint(size), img2, resize.NearestNeighbor)
	}

	newWidth := img2.Bounds().Max.X
	newHeight := img2.Bounds().Max.Y

	kittyOpts := rasterm.KittyImgOpts{SrcWidth: uint32(newWidth), SrcHeight: uint32(newHeight)}

	switch {
	case rasterm.IsKittyCapable():
		return rasterm.KittyWriteImage(os.Stdout, img2, kittyOpts)

	case rasterm.IsItermCapable():
		return rasterm.ItermWriteImage(os.Stdout, img2)

	case sixelCapable:
		// Convert image to a paletted format
		palettedImg := ConvertToPaletted(img2)
		return rasterm.SixelWriteImage(os.Stdout, palettedImg)

	default:
		// Ascii art fallback
		fmt.Print("\n", ConvertToASCII(img2)) // Align image at the initial position instead of \n first?
	}
	return nil
}

func DecodeImage(imageData []byte) (*image.Image, error) {
	imageDataReader := bytes.NewReader(imageData)
	image, _, err := image.Decode(imageDataReader)
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// ConvertToPaletted converts an image.Image to an image.Paletted
// Needed for Sixel conversion
func ConvertToPaletted(img image.Image) *image.Paletted {
	bounds := img.Bounds()
	palettedImg := image.NewPaletted(bounds, palette.Plan9)
	draw.Draw(palettedImg, bounds, img, bounds.Min, draw.Over)
	return palettedImg
}

// ASCII art conversion
func ConvertToASCII(img image.Image) string {
	converter := ascii.NewImageConverter()
	convertOptions := ascii.DefaultOptions
	return converter.Image2ASCIIString(img, &convertOptions)
}

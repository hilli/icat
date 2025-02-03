package icat

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/BourgeoisBear/rasterm"
	"github.com/hilli/icat/util"
	"github.com/nfnt/resize"

	ascii "github.com/qeesung/image2ascii/convert"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/riff"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/vp8"
	_ "golang.org/x/image/vp8l"
	_ "golang.org/x/image/webp"
)

func PrintImageFile(imageFileName string) error {
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
	return PrintImage(img, imageFileName, imageSize)
}

func PrintImageURL(imageURL string) error {
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
	return PrintImage(img, imageURL, resp.ContentLength)
}

func PrintImage(img *image.Image, filename string, imageSize int64) error {
	var img2 image.Image = *img
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
		// TODO: Convert image to a paletted format
		// if iPaletted, bOK := img.(*image.Paletted); bOK {
		// 	return rasterm.SixelWriteImage(os.Stdout, iPaletted)
		// } else {
		// 	fmt.Println("[NOT PALETTED, SKIPPING.]")
		// 	return nil
		// }

	default:
		// Ascii art fallback
		converter := ascii.NewImageConverter()
		convertOptions := ascii.DefaultOptions
		fmt.Print("\n", converter.Image2ASCIIString(img2, &convertOptions)) // Align image at the initial position instead of \n first?
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

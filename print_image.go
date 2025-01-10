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

	imageConfig, err := DecodeImageConfig(imageData)
	if err != nil {
		return err
	}

	img, err := DecodeImage(imageData)
	if err != nil {
		return err
	}
	return PrintImage(img, &imageConfig, imageFileName, imageSize)
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

	imageConfig, err := DecodeImageConfig(imageData)
	if err != nil {
		fmt.Println("Error decoding image config:", err)
		return err
	}

	img, err := DecodeImage(imageData)
	if err != nil {
		return err
	}
	return PrintImage(img, &imageConfig, imageURL, resp.ContentLength)
}

func PrintImage(img *image.Image, imageConfig *image.Config, filename string, imageSize int64) error {
	sixelCapable, _ := rasterm.IsSixelCapable()

	_, _, pw, ph := TermSize() // Get terminal height and width in pixels

	kittyOpts := rasterm.KittyImgOpts{SrcWidth: uint32(imageConfig.Width), SrcHeight: uint32(imageConfig.Height)}

	if pw < uint16(imageConfig.Width) {
		kittyOpts.SrcWidth = uint32(pw)
	}
	if ph < uint16(imageConfig.Height) {
		kittyOpts.SrcHeight = uint32(ph)
	}

	// resizedImage := resizeImage(*img, kittyOpts.SrcHeight)

	switch {
	case rasterm.IsKittyCapable():
		return rasterm.KittyWriteImage(os.Stdout, *img, kittyOpts)

	case rasterm.IsItermCapable():
		return rasterm.ItermWriteImage(os.Stdout, *img)

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
		fmt.Print("\n", converter.Image2ASCIIString(*img, &convertOptions)) // Align image at the initial position instead of \n first?
	}
	return nil
}

func DecodeImageConfig(imageData []byte) (imageConfig image.Config, err error) {
	imageDataReader := bytes.NewReader(imageData)
	imageConfig, _, err = image.DecodeConfig(imageDataReader)
	if err != nil {
		return image.Config{}, err
	}
	return imageConfig, nil
}

func DecodeImage(imageData []byte) (*image.Image, error) {
	imageDataReader := bytes.NewReader(imageData)
	image, _, err := image.Decode(imageDataReader)
	if err != nil {
		return nil, err
	}
	return &image, nil
}

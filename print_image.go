package icat

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"

	"github.com/BourgeoisBear/rasterm"
	"github.com/hilli/icat/util"
	ascii "github.com/qeesung/image2ascii/convert"
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

	fmt.Printf("image size: %d bytes\n", len(imageData))
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
	resp, err := http.Get(imageURL)
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

func PrintImage(image *image.Image, imageConfig *image.Config, filename string, imageSize int64) error {
	sixelCapable, _ := rasterm.IsSixelCapable()

	_, _, pw, ph := TermSize() // Get terminal height and width in pixels

	kittyOpts := rasterm.KittyImgOpts{SrcWidth: uint32(imageConfig.Width), SrcHeight: uint32(imageConfig.Height)}

	if pw < uint16(imageConfig.Width) {
		kittyOpts.SrcWidth = uint32(pw)
	}
	if ph < uint16(imageConfig.Height) {
		kittyOpts.SrcHeight = uint32(ph)
	}

	fmt.Println("kittyOpts:", kittyOpts)
	switch {
	case rasterm.IsKittyCapable():
		return rasterm.KittyWriteImage(os.Stdout, *image, kittyOpts)

	case rasterm.IsItermCapable():
		rasterm.ItermWriteImageWithOptions(os.Stdout, *image, rasterm.ItermImgOpts{Width: string(imageConfig.Width), Height: string(imageConfig.Height), Name: filename, DisplayInline: true})
		// rasterm.ItermCopyFileInlineWithOptions()
		return rasterm.ItermWriteImage(os.Stdout, *image)

	case sixelCapable:
		// TODO: Convert image to a paletted format
		// return rasterm.SixelWriteImage(os.Stdout, *image)

	default:
		// Ascii art fallback
		converter := ascii.NewImageConverter()
		convertOptions := ascii.DefaultOptions
		fmt.Print("\n", converter.Image2ASCIIString(*image, &convertOptions)) // Align image at the initial position instead of \n first?
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

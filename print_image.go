package icat

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"

	"github.com/BourgeoisBear/rasterm"
	"github.com/hilli/icat/util"
	"github.com/qeesung/image2ascii/convert"
	"golang.org/x/image/webp"
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
	imageConfig, imageType, err := DecodeImageConfig(imageData)
	if err != nil {
		return err
	}

	// _, err = imageFile.Seek(0, 0)
	// if err != nil {
	// 	return err
	// }

	// image, _, err := image.Decode(imageFile)
	img, err := DecodeImage(imageData, imageType)
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

	imageConfig, imageType, err := DecodeImageConfig(imageData)
	if err != nil {
		fmt.Println("Error decoding image config:", err)
		return err
	}

	fmt.Println("Image type:", imageType)

	img, err := DecodeImage(imageData, imageType)
	if err != nil {
		return err
	}
	return PrintImage(img, &imageConfig, imageURL, resp.ContentLength)
}

func PrintImage(image *image.Image, imageConfig *image.Config, filename string, imageSize int64) error {
	sixelCapable, _ := rasterm.IsSixelCapable()
	switch {
	case rasterm.IsKittyCapable():
		return rasterm.KittyWriteImage(os.Stdout, *image, rasterm.KittyImgOpts{SrcWidth: uint32(imageConfig.Width), SrcHeight: uint32(imageConfig.Height)})

	case rasterm.IsItermCapable():
		rasterm.ItermWriteImageWithOptions(os.Stdout, *image, rasterm.ItermImgOpts{Width: string(imageConfig.Width), Height: string(imageConfig.Height), Name: filename})
		// rasterm.ItermCopyFileInlineWithOptions()
		return rasterm.ItermWriteImage(os.Stdout, *image)
	case sixelCapable:
		// Convert image to a paletted format
		// return rasterm.SixelWriteImage(os.Stdout, *image)
	default:
		// Ascii art fallback
		converter := convert.NewImageConverter()
		convertOptions := convert.DefaultOptions
		fmt.Print("\n", converter.Image2ASCIIString(*image, &convertOptions)) // Align image at the initial position instead of \n first?
	}
	return nil
}

func DecodeImageConfig(imageData []byte) (imageConfig image.Config, imageType string, err error) {
	imageDataReader := bytes.NewReader(imageData)
	imageConfig, imageType, err = image.DecodeConfig(imageDataReader)

	if err != nil && errors.Is(err, image.ErrFormat) {
		// Lets try webp
		_, _ = imageDataReader.Seek(0, 0)
		imageConfig, err = webp.DecodeConfig(imageDataReader)
		fmt.Println("imageConfig:", imageConfig, err)
		if err != nil {
			return image.Config{}, "", err
		}
		imageType = "webp"
	}
	return imageConfig, imageType, nil
}

func DecodeImage(imageData []byte, imageType string) (*image.Image, error) {
	imageDataReader := bytes.NewReader(imageData)
	if imageType == "webp" {
		img, err := webp.Decode(imageDataReader)
		return &img, err
	}
	image, _, err := image.Decode(imageDataReader)
	if err != nil {
		return nil, err
	}
	return &image, nil
}

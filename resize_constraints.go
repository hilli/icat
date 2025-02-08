package icat

import (
	"image"
)

func resizeConstraints(imageBounds image.Rectangle, maxHeight, maxWidth int) (size int, option rune) {
	imgHeight := imageBounds.Dy()
	imgWidth := imageBounds.Dx()

	if maxHeight == 0 || maxWidth == 0 {
		return 0, '0' // Don't resize
	}

	if imgHeight <= maxHeight && imgWidth <= maxWidth {
		return 0, '0' // Don't resize
	}

	if imgHeight <= maxHeight && imgWidth > maxWidth {
		return maxWidth, 'x'
	}

	if imgWidth <= maxWidth && imgHeight > maxHeight {
		return maxHeight, 'y'
	}

	// OK, both x and y are too big. Let's figure out which one is the biggest
	// and use that to calculate the resize factor.
	hRatio := float32(imgHeight) / float32(maxHeight)
	wRatio := float32(imgWidth) / float32(maxWidth)
	if hRatio > wRatio {
		return maxHeight, 'y'
	} else {
		return maxWidth, 'x'
	}
}

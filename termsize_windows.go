//go:build windows
// +build windows

package icat

// TermSize returns the size of the terminal in rows and columns, as well as the pixel width and height.
// This is not supported on Windows.
// If the terminal size cannot be determined, it returns (0, 0, 0, 0).
func TermSize() (rows, columns, pixelWith, pixelHeight uint16) {
	return 0, 0, 0, 0
}

//go:build !windows
// +build !windows

package icat

import (
	"os"

	"golang.org/x/sys/unix"
)

// TermSize returns the size of the terminal in rows and columns, as well as the pixel width and height.
// The latter might not be supported on all terminals.
// If the terminal size cannot be determined, it returns (0, 0, 0, 0).
func TermSize() (rows, columns, pixelWith, pixelHeight uint16) {
	var err error
	var f *os.File
	if f, err = os.OpenFile("/dev/tty", unix.O_NOCTTY|unix.O_CLOEXEC|unix.O_NDELAY|unix.O_RDWR, 0666); err == nil {
		var sz *unix.Winsize
		if sz, err = unix.IoctlGetWinsize(int(f.Fd()), unix.TIOCGWINSZ); err == nil {
			return sz.Row, sz.Col, sz.Xpixel, sz.Ypixel
		}
	}
	return 0, 0, 0, 0
}

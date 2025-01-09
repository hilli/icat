//go:build heif

// Set --tags=heif to build support for hief in to icat.
// libheif is in LGPL-3.0 license, so it is not included in the default build.
// This also requires the libheif library to be installed on the system and using CGO.

package icat

import (
	_ "github.com/strukturag/libheif/go/heif"
)

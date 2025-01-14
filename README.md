# icat

Go version of image cat (or `imgcat`) with `bmp`, `riff`,`tiff`, `vp8`, `vp8l`, `webp` and (optional) `HEIF` image support (on top of standard Go images formats `gif`, `jpg`, `png`) and URL support and fallback to ASCII art.
`PrintImageFile` and `PrintImageURL` will print the image to the terminal, trying to figure out which terminal you are using and fallback to ASCII art if the terminal does not support images.

## Cmd line version

### Installation

```shell
go install github.com/hilli/icat/cmd/icat@latest
# With HEIF support:
go install -tags heif github.com/hilli/icat/cmd/icat@latest
```

*Note*: To enable `heif` support, you need to have `libheif` installed on your system. On macOS, you can install it via Homebrew:

```shell
brew install libheif
```

On Ubuntu, you can install it via `apt`:

```shell
sudo apt install libheif-dev
```

but you will also probably need to downgrade the version used in `icat` before compilation. Ie installing libheif on Ubuntu 24.10 will [currently] give you version 1.18.1:

```shell
go get github.com/strukturag/libheif@v1.18.1
```

and then you can install `icat` with `heif` support:

```shell
go install --tags=heif cmd/icat/icat.go
ls -l $(go env GOPATH)/bin/icat
```

`libheif` should be available on most package managers.

Alternatively, you can download the binary from the [releases page](https://github.com/hilli/icat/releases) or install via [Homebrew](https://brew.sh/):

```shell
brew install hilli/tap/icat
```

which unfortunately does not support the `heif` files (Cross compilation is hard with C extensions, mkay).

### CLI Usage

```shell
icat image.jpg
# Or
icat https://example.com/image.png
```

## Library

### Library Usage

```go
package main

import (
  "fmt"
  "os"

  "github.com/hilli/icat"
)

func main() {
  err := icat.PrintImageFile("image.jpg")
  if err != nil {
    fmt.Println(err)
  }
}
```


## LICENSE

[MIT](LICENSE)

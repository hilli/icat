# icat

Go version of image cat (or `imgcat`) with `webp`, `tiff` and (optional) `HEIF` image support (on top of standard Go images formats `gif`, `jpg`, `png`) and URL support and fallback to ASCII art.
`PrintImageFile` and `PrintImageURL` will print the image to the terminal, trying to figure out which terminal you are using and fallback to ASCII art if the terminal does not support images.

## Cmd line version

### Installation

```shell
go install github.com/hilli/icat/cmd/icat@latest
```

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

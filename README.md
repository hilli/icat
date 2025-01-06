# icat

Go version of image cat (or `imgcat`) with `webp` and URL support and fallback to ASCII art.

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

package main

import (
	"net/url"
	"os"

	"github.com/hilli/icat"
)

func main() {

	args := os.Args[1:]

	if len(args) == 0 {
		os.Exit(0)
	}

	for _, arg := range args {
		var err error

		url, err := url.Parse(arg)
		if err == nil && (url.Scheme == "http" || url.Scheme == "https") {
			err = icat.PrintImageURL(arg)
		} else {
			err = icat.PrintImageFile(arg)
		}
		if err != nil {
			os.Stderr.WriteString(err.Error())
		}
	}
}

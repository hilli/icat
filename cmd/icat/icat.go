package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/hilli/icat"
)

func main() {
	// Define the -ascii flag
	forceASCII := flag.Bool("ascii", false, "Force ASCII art output")

	// Define a custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <image-url-or-path>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Get the remaining arguments after flag parsing
	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	for i, arg := range args {
		var err error

		url, err := url.Parse(arg)
		if err == nil && (url.Scheme == "http" || url.Scheme == "https") {
			err = icat.PrintImageURL(arg, *forceASCII)
		} else {
			err = icat.PrintImageFile(arg, *forceASCII)
		}
		if err != nil {
			os.Stderr.WriteString(err.Error())
		}
		if i < len(args)-1 {
			os.Stdout.WriteString("\n")
		}
	}
}

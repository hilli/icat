package main

import (
	"flag"

	"os"

	"github.com/hilli/icat"
)

func main() {

	urlFlag := flag.Bool("u", false, "URL? Set to true if the argument is a URL")

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		os.Exit(0)
	}

	for _, arg := range args {
		var err error
		if *urlFlag {
			err = icat.PrintImageURL(arg)
		} else {
			err = icat.PrintImageFile(arg)
		}
		if err != nil {
			os.Stderr.WriteString(err.Error())
		}
	}
}

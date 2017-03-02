package main

import (
	"os"

	"github.com/martinlindhe/imgcat/lib"
)

func main() {
	inFile := "/Users/bkazi/Downloads/gotta-go-fast.gif"

	imgcat.CatFile(inFile, os.Stdout)
}

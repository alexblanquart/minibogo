package main

import (
	"flag"
	"image"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

var target string

func resizer(toWalk string, f os.FileInfo, err error) error {
	file, err := os.Open(toWalk)
	defer file.Close()

	initial, _, err := image.Decode(file)
	if err != nil {
		// Happen for root directory and images with unkown/unhandled format.
		// For now : png and jpeg!
		return nil
	}

	// resize to 340 of heigth using Lanczos resampling
	// and preserve aspect ratio
	resized := resize.Resize(340, 0, initial, resize.Lanczos3)

	// get name of image without the extra stuff to compute new path
	name := filepath.Base(toWalk)
	ext := filepath.Ext(toWalk)
	nameWithoutExt := name[:len(name)-len(ext)]
	newPath := target + nameWithoutExt + ".png"
	log.Printf("%s image about to created", newPath)

	// create the new file
	out, err := os.Create(newPath)
	if err != nil {
		log.Printf("%s when creating new file", err)
		return err
	}
	defer out.Close()

	// write new image to file
	return png.Encode(out, resized)
	if err != nil {
		log.Printf("%s when encoding new file as png image", err)
	}
	return nil
}

func main() {
	// get directory to walk : the 1st argument
	// and directory where computed images will be located : the 2nd argument
	flag.Parse()
	if flag.NArg() != 2 {
		log.Printf("usage : resizer directoryToWalk directoryWithComputedImages ")
		os.Exit(2)
	}
	toWalk := flag.Arg(0)
	target = flag.Arg(1)

	// renew thumbnails directory no matter what
	os.RemoveAll(target)
	err := os.Mkdir(target, 0777)
	if err != nil {
		log.Fatal(err)
	}

	// call resizer on each walked file
	filepath.Walk(toWalk, resizer)
}

package main

import (
	"flag"
	"github.com/nfnt/resize"
	"image"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

func resizer(path string, f os.FileInfo, err error) error {
	file, err := os.Open(path)
	defer file.Close()

	initial, _, err := image.Decode(file)
	if err != nil {
		// Happen for root directory and images with unkown/unhandled format.
		// For now : png and jpeg!
		return nil
	}

	// resize to 340 of heigth using Lanczos resampling
	// and preserve aspect ratio
	resized := resize.Resize(0, 340, initial, resize.Lanczos3)

	// get name of image without the extra stuff to compute new path
	name := filepath.Base(path)
	ext := filepath.Ext(path)
	nameWithoutExt := name[:len(name)-len(ext)]
	newPath := "static/images/thumbs/" + nameWithoutExt + ".png"
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
	// get directory to walk
	flag.Parse()
	root := flag.Arg(0) // 1st argument is the directory location

	// renew thumbnails directory no matter what
	thumbsPath := "static/images/thumbs"
	os.RemoveAll(thumbsPath)
	err := os.Mkdir(thumbsPath, 0777)
	if err != nil {
		log.Fatal(err)
	}

	// call resizer on each walked file
	filepath.Walk(root, resizer)
}

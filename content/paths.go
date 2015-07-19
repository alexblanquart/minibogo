package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// A specific function replacing every line refering to images from ![](image.jpg) to ![](/static/images/thumbs/image.png) 
// Not meant to be used more than once. Just to prepare the original content based on the present workspace.
func imagesPaths(contentPath string, info os.FileInfo, err error) error {
	if info.IsDir() || !strings.HasSuffix(contentPath, ".md") {
		// Directory and non markdown file are just ignored
		return nil
	}

	input, err := ioutil.ReadFile(contentPath)
	if err != nil {
		return err
	}

	output := bytes.Replace(input, []byte("![]("), []byte("![](/static/images/thumbs/"), -1)
	output = bytes.Replace(output, []byte(".jpg"), []byte(".png"), -1)

	if err = ioutil.WriteFile(contentPath, output, 0666); err != nil {
		return err
	}

	return nil
}

func main() {
	filepath.Walk(".", imagesPaths)
}

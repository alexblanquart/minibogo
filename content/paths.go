package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

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

package main

import (
	compression "github.com/nurlantulemisov/imagecompression"
	"image/png"
	"log"
	"os"
)

func main() {
	file, err := os.Open("examples/simple_usage/tmp/test.png")

	if err != nil {
		log.Fatalf(err.Error())
	}

	img, err := png.Decode(file)

	if err != nil {
		log.Fatalf(err.Error())
	}

	compressing, _ := compression.New(95)
	compressingImage := compressing.Compress(img)

	f, err := os.Create("examples/simple_usage/tmp/compressed-test.png")
	if err != nil {
		log.Fatalf("error creating file: %s", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf(err.Error())
		}
	}(f)

	err = png.Encode(f, compressingImage)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

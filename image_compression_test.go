package image_compression

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var testImageCompressionProvider = []struct {
	ratio       int
	inputImg    image.Image
	path        string
	expectedImg image.Image
}{
	{
		ratio:       10,
		inputImg:    openImageFromFilePath("fixtures/test10x10.jpeg"),
		path:        "fixtures/test10x10_actual10.jpeg",
		expectedImg: openImageFromFilePath("fixtures/test10x10_compressed10.jpeg"),
	},
	{
		ratio:       50,
		inputImg:    openImageFromFilePath("fixtures/test10x10.jpeg"),
		path:        "fixtures/test10x10_actual50.jpeg",
		expectedImg: openImageFromFilePath("fixtures/test10x10_compressed50.jpeg"),
	},
	{
		ratio:       80,
		inputImg:    openImageFromFilePath("fixtures/test10x10.jpeg"),
		path:        "fixtures/test10x10_actual80.jpeg",
		expectedImg: openImageFromFilePath("fixtures/test10x10_compressed80.jpeg"),
	},
	{
		ratio:       99,
		inputImg:    openImageFromFilePath("fixtures/test10x10.jpeg"),
		path:        "fixtures/test10x10_actual99.jpeg",
		expectedImg: openImageFromFilePath("fixtures/test10x10_compressed99.jpeg"),
	},
}

func TestNewWithError(t *testing.T) {
	_, err := New(100)
	if err == nil {
		t.Errorf("New expected error")
	}
}

func TestCompress(t *testing.T) {
	for _, tc := range testImageCompressionProvider {
		imgCompression, _ := New(tc.ratio)
		actual := imgCompression.Compress(tc.inputImg)
		saveImage(actual, tc.path)

		actualImage := openImageFromFilePath(tc.path)
		if !reflect.DeepEqual(tc.expectedImg, actualImage) {
			t.Errorf("Matrix %v expected, actual %v", tc.expectedImg, actual)
		}
		removeImage(tc.path)
	}
}

func openImageFromFilePath(filePath string) image.Image {
	f, _ := os.Open(filePath)
	defer f.Close()

	var img image.Image
	ext := filepath.Ext(filePath)
	switch ext {
	case ".png":
		img, _ = png.Decode(f)
	case ".jpeg", ".jpg":
		img, _ = jpeg.Decode(f)
	default:
		panic(fmt.Errorf("unsupported format %s", ext))
	}

	return img
}

func saveImage(img image.Image, path string) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("error creating file: %s", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("error creating file: %s", err)
		}
	}(f)

	err = jpeg.Encode(f, img, nil)
	if err != nil {
		return
	}
}

func removeImage(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Fatalf("error creating file: %s", err)
		return
	}
}

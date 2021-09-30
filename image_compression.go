package image_compression

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"image"
	"image/color"
	"math"
)

// ImageCompression provide service for compression image.Image
type ImageCompression struct {
	// ratio is percent compression.
	// 0 <= ratio < 100, 0 is original images
	ratio int
}

func New(ratio int) (*ImageCompression, error) {
	if ratio >= 100 {
		return nil, fmt.Errorf("ratio can't be more or equal 100, got %d", ratio)
	}
	return &ImageCompression{ratio: ratio}, nil
}

// Compress return new image.Image with ratio
func (im *ImageCompression) Compress(img image.Image) image.Image {
	if im.ratio == 0 {
		return img
	}

	width, height := img.Bounds().Max.X, img.Bounds().Max.Y
	rank := im.toMode(width, height)

	var redData, greenData, blueData, alphaData []float64
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			redData = append(redData, float64(r>>8))
			greenData = append(greenData, float64(g>>8))
			blueData = append(blueData, float64(b>>8))
			alphaData = append(alphaData, float64(a>>8))
		}
	}

	return im.compressingImage(width, height, rank, &redData, &greenData, &blueData, &alphaData)
}

// compressingImage is sub-function which decomposing RGBA slices and approximate by rank.
// return compressed image.Image
func (im *ImageCompression) compressingImage(width, height, rank int, redData, greenData, blueData, alphaData *[]float64) image.Image {
	chRed, chGreen, chBlue, chAlpha := make(chan mat.Dense), make(chan mat.Dense), make(chan mat.Dense), make(chan mat.Dense)

	go im.approximateImgChannel(width, height, rank, redData, chRed)
	go im.approximateImgChannel(width, height, rank, greenData, chGreen)
	go im.approximateImgChannel(width, height, rank, blueData, chBlue)
	go im.approximateImgChannel(width, height, rank, alphaData, chAlpha)

	var compressRedChannel, compressGreenChannel, compressBlueChannel, compressAlphaChannel mat.Dense

	compressRedChannel = <-chRed
	compressGreenChannel = <-chGreen
	compressBlueChannel = <-chBlue
	compressAlphaChannel = <-chAlpha

	var newImg = image.NewRGBA(image.Rect(0, 0, width, height))
	var col color.Color

	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			col = color.RGBA{
				R: uint8(compressRedChannel.At(r, c)),
				G: uint8(compressGreenChannel.At(r, c)),
				B: uint8(compressBlueChannel.At(r, c)),
				A: uint8(compressAlphaChannel.At(r, c)),
			}
			newImg.Set(c, r, col)
		}
	}
	return newImg
}

func (im *ImageCompression) approximateImgChannel(width, height, rank int, imgChannel *[]float64, ch chan mat.Dense) {
	ch <- approximate(mat.NewDense(height, width, *imgChannel), rank)
}

// toMode convert ratio percent to k elements from singular values
func (im *ImageCompression) toMode(width, height int) int {
	return int(math.Ceil((float64(min(width, height)) / 100.0) * float64(100-im.ratio)))
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

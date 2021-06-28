package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"net/http"
	"strconv"
	"syscall/js"

	"github.com/anthonynsimon/bild/imgio"
	// "github.com/anthonynsimon/bild/transform"
	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/convolution"
	"github.com/anthonynsimon/bild/parallel"
)

func main() {

	
	resp, err := http.Get("lenna.png")
	fmt.Println(err)
	
	//type- OriginalImgRGBA =  *image.RGBA
	OriginalImgRGBA, err := png.Decode(resp.Body)
	if err != nil {
		log.Fatalf("Error decoding PNG: %s", err.Error())
	}

	// Create grayscale of img
	newImgRGBA := Gaussian(OriginalImgRGBA, 9.0)
	//fmt.Println(newImgRGBA)

	buf := new(bytes.Buffer)
	encoder := imgio.PNGEncoder()
	encoder(buf, newImgRGBA)
	newBitmap := buf.Bytes()
	//fmt.Println("buf bytes: ", newBitmap)

	// Type of new bitmap in Go- []uint8
	// fmt.Println(reflect.TypeOf(newBitmap));

	console := js.Global().Get("console")

	dst := js.Global().Get("Uint8Array").New(len(newBitmap))
	n := js.CopyBytesToJS(dst, newBitmap)
	console.Call("log", "bytes copied:", strconv.Itoa(n))
	js.Global().Call("displayImage", dst)

}

func Grayscale(img image.Image) *image.RGBA {
	return GrayscaleWithWeights(img, 0.3, 0.6, 0.1)
}

func GrayscaleWithWeights(img image.Image, r, g, b float64) *image.RGBA {
	src := clone.AsRGBA(img)
	bounds := src.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()

	if bounds.Empty() {
		return &image.RGBA{}
	}

	dst := image.NewRGBA(bounds)

	parallel.Line(srcH, func(start, end int) {
		for y := start; y < end; y++ {
			for x := 0; x < srcW; x++ {
				pos := y*src.Stride + x*4

				c := r*float64(src.Pix[pos+0]) + g*float64(src.Pix[pos+1]) + b*float64(src.Pix[pos+2])
				k := uint8(c + 0.5)
				dst.Pix[pos] = k
				dst.Pix[pos+1] = k
				dst.Pix[pos+2] = k
				dst.Pix[pos+3] = src.Pix[pos+3]
			}
		}
	})
	return dst
}
	
func Gaussian(src image.Image, radius float64) *image.RGBA {
	if radius <= 0 {
		return clone.AsRGBA(src)
	}

	// Create the 1-d gaussian kernel
	length := int(math.Ceil(2*radius + 1))
	k := convolution.NewKernel(length, 1)
	for i, x := 0, -radius; i < length; i, x = i+1, x+1 {
		k.Matrix[i] = math.Exp(-(x * x / 4 / radius))
	}
	normK := k.Normalized()

	// Perform separable convolution
	options := convolution.Options{Bias: 0, Wrap: false, KeepAlpha: false}
	result := convolution.Convolve(src, normK, &options)
	result = convolution.Convolve(result, normK.Transposed(), &options)

	return result
}

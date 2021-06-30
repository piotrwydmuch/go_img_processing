package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"net/http"
	"strconv"
	"sync"
	"syscall/js"
	"time"

	"github.com/anthonynsimon/bild/imgio"
	// "github.com/anthonynsimon/bild/transform"
	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/convolution"
	"github.com/anthonynsimon/bild/parallel"
)

func main() {

	//Loading and decoding sample image
	resp, err := http.Get("lenna.png")
	fmt.Println(err)	
	OriginalImgRGBA, err := png.Decode(resp.Body)
	if err != nil {
		log.Fatalf("Error decoding PNG: %s", err.Error())
	}
	
	//JavaScript func declaratrion - calling from DOM is possible now
	var gaussianFunc js.Func
	gaussianFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		UseGaussian(OriginalImgRGBA)
		//gaussianFunc.Release() // release the function if the button will not be clicked again
		return nil
	})
	js.Global().Get("document").Call("getElementById", "btn-1").Call("addEventListener", "click", gaussianFunc)
	
	var grayscaleFunc js.Func
	grayscaleFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		UseGrayscale(OriginalImgRGBA)
		//cb.Release() // release the function if the button will not be clicked again
		return nil
	})
	js.Global().Get("document").Call("getElementById", "btn-2").Call("addEventListener", "click", grayscaleFunc)

	
	var invertFunc js.Func
	invertFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		UseInvert(OriginalImgRGBA)
		//cb.Release() // release the function if the button will not be clicked again
		return nil
	})
	js.Global().Get("document").Call("getElementById", "btn-3").Call("addEventListener", "click", invertFunc)
	
		
	//App is ready for user actions
	fmt.Println("App ready.")	
	
	//To wait for multiple goroutines to finish, we can use a wait group
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func UseGaussian(img image.Image) {
	newImgRGBA := Gaussian(img, 9.0)
	EditPhoto(newImgRGBA)
}

func UseGrayscale(img image.Image) {
	newImgRGBA := Grayscale(img)
	EditPhoto(newImgRGBA)
}

func UseInvert(img image.Image) {
	newImgRGBA := Invert(img)
	EditPhoto(newImgRGBA)
}


func EditPhoto(img image.Image) {
	//Measuring execution time 
	start := time.Now()
	
	//bitmap create
	buf := new(bytes.Buffer)
	encoder := imgio.PNGEncoder()
	encoder(buf, img)
	newBitmap := buf.Bytes()
	
	//get the browser console object
	console := js.Global().Get("console")
	
	//bitman -> base64 result needed for representation of processed image
	dst := js.Global().Get("Uint8Array").New(len(newBitmap))
	n := js.CopyBytesToJS(dst, newBitmap)
	console.Call("log", "bytes copied:", strconv.Itoa(n))
	js.Global().Call("displayImage", dst)
	
	fmt.Println("t: ", time.Since(start))
}

func Invert(src image.Image) *image.RGBA {
	fn := func(c color.RGBA) color.RGBA {
		return color.RGBA{255 - c.R, 255 - c.G, 255 - c.B, c.A}
	}

	img := adjust.Apply(src, fn)

	return img
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

func Grayscale(img image.Image) *image.RGBA {
	return GrayscaleWithWeights(img, 0.3, 0.6, 0.1)
}


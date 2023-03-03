//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"log"
	"math"
	"syscall/js"
)

// taken from stdlib's image/image.go (SubImage) and modified to account for
// gif frame positions
func Crop(p *image.Paletted, selection image.Rectangle) *image.Paletted {
	intersection := p.Rect.Intersect(selection)

	if intersection.Empty() {
		return &image.Paletted{
			Palette: p.Palette,
		}
	}

	i := p.PixOffset(intersection.Min.X, intersection.Min.Y)
	xMin := int(math.Round(float64(intersection.Min.X) - (float64(selection.Min.X))))
	yMin := int(math.Round(float64(intersection.Min.Y) - (float64(selection.Min.Y))))
	xMax := xMin + (intersection.Dx())
	yMax := yMin + (intersection.Dy())
	bounds := image.Rect(xMin, yMin, xMax, yMax)

	return &image.Paletted{
		Pix:     p.Pix[i:],
		Stride:  p.Stride,
		Rect:    bounds,
		Palette: p.Palette,
	}
}

func encodeGif(this js.Value, args []js.Value) interface{} {
	rawGifBytes := args[0]
	x := args[1].Int()
	y := args[2].Int()
	w := args[3].Int()
	h := args[4].Int()

	gifBytes := make([]byte, rawGifBytes.Length())
	js.CopyBytesToGo(gifBytes, rawGifBytes)
	gifIn, err := gif.DecodeAll(bytes.NewReader(gifBytes))
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(gifIn.Image); i++ {
		frame := gifIn.Image[i]
		selection := image.Rect(x, y, x+w, y+h)
		gifIn.Image[i] = Crop(frame, selection)
	}

	outGif := gif.GIF{
		Image:           gifIn.Image,
		Delay:           gifIn.Delay,
		LoopCount:       gifIn.LoopCount,
		Disposal:        gifIn.Disposal,
		BackgroundIndex: gifIn.BackgroundIndex,
	}

	outBytes := new(bytes.Buffer)
	err = gif.EncodeAll(outBytes, &outGif)
	if err != nil {
		log.Fatal(err)
	}

	return js.CopyBytesToJS(args[5], outBytes.Bytes())
}

func main() {
	done := make(chan struct{}, 0)
	js.Global().Set("encodeGif", js.FuncOf(encodeGif))
	<-done
}

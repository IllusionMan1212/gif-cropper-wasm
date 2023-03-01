//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"log"
	"syscall/js"
)

func encodeGif(this js.Value, args []js.Value) interface{} {
	inImages := args[0]
	delays := args[1]
	disposals := args[2]

	outImages := make([]*image.Paletted, 0)
	delay := make([]int, 0)
	disposal := make([]byte, 0)

	// Plan9 palette doesn't include transparent colors
	// so we use WebSafe and add transparent color
	myPalette := append(palette.WebSafe, image.Transparent)
	opts := gif.Options{
		NumColors: 256,
		Drawer:    draw.FloydSteinberg,
	}

	fmt.Println(inImages.Length())
	for i := 0; i < inImages.Length(); i++ {
		imgBytes := make([]byte, inImages.Index(i).Length())
		js.CopyBytesToGo(imgBytes, inImages.Index(i))

		reader := bytes.NewReader(imgBytes)
		img, err := png.Decode(reader)
		if err != nil {
			log.Fatal(err)
		}

		b := img.Bounds()
		pimg := image.NewPaletted(b, myPalette)
		opts.Drawer.Draw(pimg, b, img, image.ZP)

		outImages = append(outImages, pimg)
		delay = append(delay, delays.Index(i).Int())
		disposal = append(disposal, byte(disposals.Index(i).Int()))
	}

	outGif := &gif.GIF{
		Image:           outImages,
		Delay:           delay,
		LoopCount:       0,
		Disposal:        disposal,
		BackgroundIndex: 0,
	}

	outBytes := new(bytes.Buffer)
	err := gif.EncodeAll(outBytes, outGif)
	if err != nil {
		log.Fatal(err)
	}

	return js.CopyBytesToJS(args[3], outBytes.Bytes())
}

func main() {
	done := make(chan struct{}, 0)
	js.Global().Set("encodeGif", js.FuncOf(encodeGif))
	<-done
}

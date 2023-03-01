package main

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No directory or file provided.")
	}

	path := os.Args[1]

	if path == "" {
		log.Fatal("No directory or file provided.")
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	images := make([]*image.Paletted, 0)
	delay := make([]int, 0)
	disposal := make([]byte, 0)

	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			continue
		}

		file, err := os.OpenFile(filepath.Join(path, fileInfo.Name()), os.O_RDWR, os.ModeExclusive)
		if err != nil {
			log.Fatal("Failed to open file")
		}
		defer file.Close()

		img, err := png.Decode(file)
		if err != nil {
			log.Fatal(err)
		}

		opts := gif.Options{
			NumColors: 256,
			Drawer:    draw.FloydSteinberg,
		}
		b := img.Bounds()
		myPalette := append(palette.WebSafe, image.Transparent)
		pimg := image.NewPaletted(b, myPalette)
		opts.Drawer.Draw(pimg, b, img, image.ZP)

		images = append(images, pimg)
		delay = append(delay, 10)
		disposal = append(disposal, 2)
	}

	outGif := &gif.GIF{
		Image:           images,
		Delay:           delay,
		LoopCount:       0,
		Disposal:        disposal,
		BackgroundIndex: 40,
	}

	out, err := os.Create("out.gif")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	err = gif.EncodeAll(out, outGif)
	if err != nil {
		log.Fatal(err)
	}
}

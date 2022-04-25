package lovetik

import (
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
)

func ApplyWatermark(ppic image.Image) (image.Image, error) {
	crownfile, err := os.Open("asset/crown.png")
	if err != nil {
		log.Fatalf("failed to open: %s", err)
	}

	second, err := png.Decode(crownfile)
	if err != nil {
		log.Fatalf("failed to decode: %s", err)
	}
	defer crownfile.Close()

	offset := image.Pt(300, 200) // AK TODO fix offset
	b := ppic.Bounds()
	image3 := image.NewRGBA(b)
	draw.Draw(image3, b, ppic, image.ZP, draw.Src)
	draw.Draw(image3, second.Bounds().Add(offset), second, image.ZP, draw.Over)

	// AK TODO return an image
	//jpeg.Encode(third, image3, &jpeg.Options{jpeg.DefaultQuality})
	return image3, nil
}

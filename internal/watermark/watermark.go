package watermark

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"time"

	log "github.com/sirupsen/logrus"
)

// Watermark operates with bytes
func Apply(bakground []byte, foreground []byte) ([]byte, error) {
	defer logDuration(time.Now())

	bgImg, _, err := image.Decode(bytes.NewReader(bakground))
	if err != nil {
		return nil, err
	}

	fgImg, err := png.Decode(bytes.NewReader(foreground))

	if err != nil {
		return nil, err
	}

	b := bgImg.Bounds()
	resImg := image.NewRGBA(b)
	draw.Draw(resImg, b, bgImg, image.Point{}, draw.Src)
	draw.Draw(resImg, fgImg.Bounds(), fgImg, image.Point{}, draw.Over)

	resBuf := new(bytes.Buffer)
	err = png.Encode(resBuf, resImg)

	if err != nil {
		return nil, err
	}

	return resBuf.Bytes(), nil
}

func logDuration(invocation time.Time) {
	elapsed := time.Since(invocation)

	log.Printf("%s lasted %s", "watermarking", elapsed)
}

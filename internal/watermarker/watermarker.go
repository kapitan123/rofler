package watermarker

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
)

type Watermarker struct {
}

func New() *Watermarker {
	return &Watermarker{}
}

func (wm *Watermarker) Apply(bg []byte, fg []byte, w io.Writer) error {
	defer logDuration(time.Now())

	bgImg, _, err := image.Decode(bytes.NewReader(bg))
	if err != nil {
		return err
	}

	fgImg, err := png.Decode(bytes.NewReader(fg))

	if err != nil {
		return err
	}

	b := bgImg.Bounds()
	resImg := image.NewRGBA(b)
	draw.Draw(resImg, b, bgImg, image.Point{}, draw.Src)
	draw.Draw(resImg, fgImg.Bounds(), fgImg, image.Point{}, draw.Over)

	err = png.Encode(w, resImg)

	if err != nil {
		return err
	}

	return nil
}

func logDuration(invocation time.Time) {
	elapsed := time.Since(invocation)

	log.Printf("%s lasted %s", "watermarking", elapsed)
}

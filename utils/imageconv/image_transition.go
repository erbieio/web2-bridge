package imageconv

import (
	"encoding/base64"
	"image"
	"image/color"

	"golang.org/x/image/draw"

	"image/gif"

	_ "embed"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/noelyahan/impexp"
	"github.com/noelyahan/mergi"
	"github.com/noelyahan/mergitrans"
)

var MaskBlack = color.RGBA{0, 0, 0, 0}
var MaskWhite = color.RGBA{255, 255, 255, 0}

//go:embed logo.png
var logoBytes []byte

//go:embed 微软雅黑.ttf
var fontBytes []byte

var logo image.Image
var fontType *truetype.Font

func init() {
	var err error
	logoBase64 := "data:image/png;base64," + base64.RawStdEncoding.EncodeToString(logoBytes)
	logo, err = mergi.Import(impexp.NewBase64Importer(logoBase64))
	if err != nil {
		panic(err)
	}
	fontType, err = truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}
}

func InkEffect(url string) (gif.GIF, error) {
	scale := 2
	logoScale := 5
	//logo, _ = mergi.Resize(logo, uint(logo.Bounds().Max.X/scale), uint(logo.Bounds().Max.Y/scale))

	img, err := mergi.Import(impexp.NewURLImporter(url))
	if err != nil {
		return gif.GIF{}, err
	}
	img, err = mergi.Resize(img, uint(img.Bounds().Max.X/scale), uint(img.Bounds().Max.Y/scale))
	if err != nil {
		return gif.GIF{}, err
	}
	resizeWidth := img.Bounds().Max.X / logoScale
	resizeHeight := logo.Bounds().Max.Y * resizeWidth / logo.Bounds().Max.X
	resizeLog := image.NewRGBA(image.Rect(0, 0, img.Bounds().Max.X/logoScale, resizeHeight))
	draw.NearestNeighbor.Scale(resizeLog, resizeLog.Rect, logo, logo.Bounds(), draw.Over, nil)

	background := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{
			X: img.Bounds().Max.X,
			Y: img.Bounds().Max.Y,
		},
	})
	draw.Draw(background, background.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	if img.Bounds().Max.X > resizeWidth && img.Bounds().Max.Y > resizeHeight {
		startX := (img.Bounds().Max.X - resizeWidth) / 2
		startY := (img.Bounds().Max.Y - resizeHeight) / 2

		dc := gg.NewContext(img.Bounds().Max.X, img.Bounds().Max.Y)
		fontPt := resizeLog.Bounds().Max.Y / 2
		dc.SetRGB(0, 0, 0)
		dc.SetFontFace(truetype.NewFace(fontType, &truetype.Options{
			Size: float64(fontPt),
		}))
		s := "Erbie"
		sWidth, sHeight := dc.MeasureString(s)
		if img.Bounds().Max.X > resizeWidth+int(sWidth) {
			startX = startX - int(sWidth/2)
			draw.Draw(background, background.Bounds(), resizeLog, image.Pt(-startX, -startY), draw.Over)

			dc.DrawString(s, float64(startX+resizeLog.Bounds().Max.X), float64(startY+resizeHeight)-(float64(resizeHeight)-sHeight)/2)
			draw.Draw(background, background.Bounds(), dc.Image(), image.Pt(0, 0), draw.Over)
		} else {
			draw.Draw(background, background.Bounds(), resizeLog, image.Pt(-startX, -startY), draw.Over)
		}

	}

	trans := mergitrans.Ink2()

	frames := mergi.Transit([]image.Image{background}, []image.Image{img}, trans, MaskBlack, 0, float64(len(trans)-1), 1)

	return mergi.Animate(frames, 1)

}

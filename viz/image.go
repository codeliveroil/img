// Copyright (c) 2018 codeliveroil. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package viz

import (
	"image"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"

	"github.com/codeliveroil/img/terminal"
	"github.com/nfnt/resize"
)

// Image is a representation of a (multi) picture
// image.
type Image struct {
	// Path to image file.
	Filename string
	// Specify a file name to export the image to a shell script.
	// For instance, this script can be used to display an image for motd.
	ExportFilename string
	// Specify a loop count to animate GIFs more than once or set to 0 to render the first picture only.
	LoopCount int
	//Specify a decimal point multiplier to increase or decrease the speed of the GIF.
	DelayMultiplier float64
	// Use specified width instead of automatically computing it. Height will be calculated according to the aspect ratio.
	// This is useful in SSH sessions where screen resizes are not registered automatically.
	UserWidth int

	frames []frame
	h      int
	w      int
}

type frame struct {
	picture [][]uint8
	delay   int
}

// Init initializes the visualization framework
// for drawing the image.
func (img *Image) Init() (err error) {
	//Open image
	file, err := os.Open(img.Filename)
	if err != nil {
		return err
	}
	firstFrame, imgFmt, err := image.Decode(file)
	if err != nil {
		return err
	}
	file.Close()

	//Identify scale
	iw := firstFrame.Bounds().Max.X
	ih := firstFrame.Bounds().Max.Y

	scale := 1.0
	if img.UserWidth > 0 {
		scale = float64(img.UserWidth) / float64(iw)
	} else {
		tw, th, err := terminal.Size()
		if err != nil {
			return err
		}
		if imgFmt == "gif" && img.LoopCount > 0 {
			tw = 40
		}
		th = (th * 2) - 1       //-1 to account for the terminal prompt ($/#) that'll show up after the image is displayed
		if tw < iw || th < ih { //scale down the image to fit the terminal
			scaleW := float64(tw) / float64(iw)
			scaleH := float64(th) / float64(ih)
			scale = math.Min(scaleW, scaleH)
		}
	}

	img.w = int(math.Floor(scale * float64(iw)))
	img.h = int(math.Floor(scale * float64(ih)))
	if img.h%2 != 0 { //make the height even since we will be painting y and y+1 in every iteration
		img.h -= 1
	}

	//Scale image frames
	appendImg := func(f image.Image, delayMS int) {
		scaled := resize.Resize(uint(img.w), uint(img.h), f, resize.Lanczos3)
		pic := make([][]uint8, img.w)
		for x := 0; x < img.w; x++ {
			pic[x] = make([]uint8, img.h)
			for y := 0; y < img.h; y++ {
				clr := scaled.At(x, y)
				x256Clr := Colors.Index(clr)
				pic[x][y] = uint8(x256Clr)
			}
		}

		img.frames = append(img.frames, frame{
			picture: pic,
			delay:   int(math.Ceil(float64(delayMS) * img.DelayMultiplier)), //GIFs will take long to render, so reduce the delay to achieve intended delay.
		})
	}

	if imgFmt == "gif" && img.LoopCount > 0 {
		file, err := os.Open(img.Filename)
		if err != nil {
			return err
		}
		g, err := gif.DecodeAll(file)
		if err != nil {
			return err
		}
		iw = g.Config.Width
		ih = g.Config.Height

		var prev *image.RGBA
		canvas := image.NewRGBA(image.Rect(0, 0, iw, ih))
		for i, frame := range g.Image {
			draw.Draw(canvas, canvas.Bounds(), frame, image.ZP, draw.Over)
			appendImg(canvas, g.Delay[i]*10)
			switch g.Disposal[i] {
			case gif.DisposalBackground:
				canvas = image.NewRGBA(image.Rect(0, 0, iw, ih))
				fallthrough
			case gif.DisposalNone:
				prev = &(*canvas)
			case gif.DisposalPrevious:
				if prev != nil {
					canvas = prev
				}
			}
		}
		file.Close()
	} else {
		img.LoopCount = 1 //override incorrect user input for single picture images
		appendImg(firstFrame, 0)
	}

	return nil
}

// Draw renders the image into one of the
// selected modes (stdout or file)
func (img *Image) Draw(canvas Canvas) error {
	firstFrameDone := false
	delay := 0
	for i := 0; i < img.LoopCount; i++ {
		for _, frame := range img.frames {
			if firstFrameDone {
				if err := canvas.LineUp(img.h / 2); err != nil {
					return err
				}
				if err := canvas.Sleep(delay); err != nil {
					return err
				}
			}
			for y := 0; y < img.h; y = y + 2 {
				for x := 0; x < img.w; x++ {
					canvas.Paint(frame.picture[x][y], frame.picture[x][y+1])
				}
				err := canvas.NewLine()
				if err != nil {
					return err
				}
			}
			firstFrameDone = true
			delay = frame.delay
		}
	}
	return canvas.Close()
}

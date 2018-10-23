// Copyright (c) 2018 codeliveroil. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package main

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/codeliveroil/img/viz"
	"github.com/codeliveroil/niceflags"
)

func main() {
	args := os.Args
	flags := niceflags.NewFlags(
		args[0],
		"Image viewer for Linux terminal emulators",
		"Supports PNG, JPEG and GIF.\n"+
			"Images can be rendered on screen (default) or exported to a shell script to be "+
			"rendered later (e.g. to display a logo during SSH login).\n"+
			"GIFs are animated and restricted to a 40 character width by default.\n"+
			"To obtain best quality rendering, try reducing the font size of the terminal.",
		"[options] file",
		"help",
		false,
	)
	flags.Examples = []string{
		"car.png",
		"logo.gif",
		"-l 2 wheel.gif",
	}

	userWidth := flags.Int("w", 0, "Use specified `width` instead of auto-computing it.")
	exportFilename := flags.String("o", "", "Export image as a shell script to specified `file`.")
	loopCount := flags.Int("l", 1, "Specify the `num`ber of times the GIF should be looped or set to 0 to render the first frame only.")
	delayMultiplier := flags.Float64("s", 1.0, "Specify a multiplier to change the `speed` of animation. "+
		"Larger the multiplier, slower the speed of animation. "+
		"For example, 2 decreases the speed to 50% and 0.5 increases the speed to 200%.")
	version := flags.Bool("v", false, "Display version.")

	check(flags.Parse(args[1:]))
	flags.Help()

	argc := len(args)
	if argc < 2 {
		flags.Usage()
		os.Exit(0)
	}
	if *version {
		fmt.Println("1.1")
		os.Exit(0)
	}
	filename := args[argc-1]
	if filename == "" {
		niceflags.PrintErr("image file not specified.\n")
		flags.Usage()
		os.Exit(1)
	}

	//Render/Export image
	img := viz.Image{
		Filename:        filename,
		ExportFilename:  *exportFilename,
		LoopCount:       *loopCount,
		DelayMultiplier: *delayMultiplier,
		UserWidth:       *userWidth,
	}

	check(img.Init())

	var canvas viz.Canvas
	if img.ExportFilename == "" {
		canvas = &viz.StdoutCanvas{}
	} else {
		var err error
		canvas, err = viz.NewFileCanvas(img.ExportFilename)
		check(err)
	}

	check(img.Draw(canvas))
}

// check prints the error message and exits
// if err is not nil.
func check(err error) {
	if err != nil {
		niceflags.PrintErr("%v\n", err)
		os.Exit(1)
	}
}

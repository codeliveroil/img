// Copyright (c) 2018 codeliveroil. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package main

import (
	"flag"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path"

	"github.com/codeliveroil/img/viz"
)

func main() {
	args := os.Args
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	print := func(msg string, args ...interface{}) {
		fmt.Fprintf(os.Stderr, msg, args...) //stderr because that's where 'flags.PrintDefaults()' prints to
	}

	fullUsage := func() {
		print("Image viewer for Linux terminal emulators.\n")
		print("  Supports PNG, JPEG and GIF.\n")
		print("  Images can be rendered on screen (default) or exported to a shell script to\n")
		print("  be rendered later (e.g. to display a logo during SSH login).\n")
		print("  GIFs are animated and restricted to a 40 character width by default.\n")
		print("  To obtain best quality rendering, try reducing the font size of the terminal.\n")

		flags.Usage()

		print("\nExamples:\n")
		print("  %s car.png\n", path.Base(os.Args[0]))
		print("  %s logo.gif\n", path.Base(os.Args[0]))
		print("  %s -l 2 wheel.gif\n", path.Base(os.Args[0]))
	}

	flags.Usage = func() {
		print("\nUsage: %s [options] file\n", path.Base(os.Args[0]))

		print("\nOptions:\n")
		printFlag := func(f *flag.Flag) { //go's default 2 line display for non-bools is not appealing, use custom printer.
			print("  -%v %v\n", f.Name, f.Usage)
		}
		flags.VisitAll(printFlag)
	}

	userWidth := flags.Int("w", 0, "num   "+
		"Use specified width instead of auto-computing it.")
	exportFilename := flags.String("o", "", "file  "+
		"Export image as a shell script to specified file.")
	loopCount := flags.Int("l", 1, "num   "+
		"Specify how many times the GIF should be looped or set to 0 to\n           "+
		"render the first frame only.")
	delayMultiplier := flags.Float64("d", 1.0, "num   "+
		"Specify a multiplier to change the speed of animation.")
	help := flags.Bool("h", false, "      "+
		"Display help screen.")
	version := flags.Bool("v", false, "      "+
		"Display version.")

	check(flags.Parse(args[1:]))
	argc := len(args)
	if *help || argc < 2 {
		fullUsage()
		os.Exit(0)
	}
	if *version {
		fmt.Println("1.0")
		os.Exit(0)
	}
	filename := args[argc-1]
	if filename == "" {
		print("image file not specified.\n")
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

	var writer viz.Writer
	if img.ExportFilename == "" {
		writer = &viz.StdWriter{}
	} else {
		var err error
		writer, err = viz.NewFileWriter(img.ExportFilename)
		check(err)
	}

	check(img.Draw(writer))
}

// check prints the error message and exits
// if err is not nil.
func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

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

	"github.com/codeliveroil/img/util"
	"github.com/codeliveroil/img/viz"
)

func main() {
	//Parse command line args
	args := os.Args
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.Usage = func() {
		print := func(msg string, args ...interface{}) {
			fmt.Fprintf(os.Stderr, msg, args...) //stderr because that's where 'flags.PrintDefaults()' prints to
		}
		print("Image viewer for Linux terminal emulators.\n")
		print("  Supports PNG, JPEG and GIF.\n")
		print("  Images can be written to a shell script to be rendered later, for instance to\n  display a company logo to a user during login.\n")
		print("  GIFs are animated and restricted to a 40 character width by default.\n")
		print("  To obtain best quality rendering, try reducing the font size of the terminal.\n")
		print("\nUsage: %s [options] file\n", path.Base(os.Args[0]))
		print("\nOptions:\n")
		flags.PrintDefaults()
		print("\nExamples:\n")
		print("%s car.png\n", path.Base(os.Args[0]))
		print("%s logo.gif\n", path.Base(os.Args[0]))
		print("%s -l 2 wheel.gif\n", path.Base(os.Args[0]))

	}

	userWidth := flags.Int("w", 0, "<width> Use specified width. Height will be calculated\n\t according to the aspect ratio. "+
		"This is useful in SSH\n\t sessions where screen resizes are not registered automatically.")
	exportFilename := flags.String("o", "", "<file_name> Export the image to a shell script. "+
		"For instance,\n this script can be used to display an image for motd.")
	loopCount := flags.Int("l", 1, "<loop_count> Specify a loop count to animate GIFs more than \n\tonce or set to 0 to render the first frame only.")
	delayMultiplier := flags.Float64("d", 1.0, "<delay_multiplier> Specify a decimal point multiplier \n\tto increase or decrease the speed of animation.")

	help := flags.Bool("h", false, "Help screen.")
	version := flags.Bool("v", false, "Show version.")

	util.Check(flags.Parse(args[1:]))
	argc := len(args)
	if *help || argc < 2 {
		flags.Usage()
		os.Exit(0)
	}
	if *version {
		fmt.Println("1.0")
		os.Exit(0)
	}
	filename := args[argc-1]
	if filename == "" {
		flags.Usage()
		os.Exit(1)
	}

	img := viz.Image{
		Filename:        filename,
		ExportFilename:  *exportFilename,
		LoopCount:       *loopCount,
		DelayMultiplier: *delayMultiplier,
		UserWidth:       *userWidth,
	}

	util.Check(img.Init())

	var writer viz.Writer
	if img.ExportFilename == "" {
		writer = &viz.StdWriter{}
	} else {
		var err error
		writer, err = viz.NewFileWriter(img.ExportFilename)
		util.Check(err)
	}

	util.Check(img.Draw(writer))
}

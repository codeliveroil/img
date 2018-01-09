// Copyright (c) 2018 codeliveroil. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/codeliveroil/img/viz"
)

func check(err error, t *testing.T) {
	if err != nil {
		t.Error("expecting no error, got", err)
	}
}

func read(filename string, t *testing.T) string {
	bytes, err := ioutil.ReadFile(filename)
	check(err, t)
	return string(bytes)
}

func export(testfile string, loopCount int, delayMultiplier float64, width int, t *testing.T) viz.Image {
	img := viz.Image{
		Filename:        testfile,
		ExportFilename:  "/tmp/img_test.sh",
		LoopCount:       loopCount,
		DelayMultiplier: delayMultiplier,
		UserWidth:       width,
	}

	os.Args = []string{"img", "-o", img.ExportFilename,
		"-l", fmt.Sprintf("%v", img.LoopCount),
		"-d", fmt.Sprintf("%v", img.DelayMultiplier),
		"-w", fmt.Sprintf("%v", img.UserWidth),
		testfile,
	}
	main() //invoke main to test flag parsing as well.

	return img
}

func validate(expected string, got viz.Image, t *testing.T) {
	if read(expected, t) != read(got.ExportFilename, t) {
		t.Errorf("expected: %v, got: %v; params: loopCount=%v, delayMultiplier=%v, userWidth=%v",
			expected, got.ExportFilename, got.LoopCount, got.DelayMultiplier, got.UserWidth)
	}
}

func TestStaticImage(t *testing.T) {
	img := export("testdata/color_matrix.png", 1, 1.0, 80, t)
	validate("testdata/color_matrix.sh", img, t)
}

func TestGIF(t *testing.T) {
	// Test different GIF disposals.
	for _, d := range []string{"Unspecified", "None", "NoneTransparency", "Background"} {
		img := export("testdata/disposal"+d+".gif", 1, 1.0, 0, t)
		validate("testdata/disposal"+d+".sh", img, t)
	}

	// Test all parameters
	img := export("testdata/disposalNone.gif", 3, 10, 60, t)
	validate("testdata/all.sh", img, t)
}

// Copyright (c) 2018 codeliveroil. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"testing"

	"github.com/codeliveroil/img/terminal"
	"github.com/codeliveroil/img/viz"
)

const testData = "resources/testData/"

func read(filename string, t *testing.T) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error("expecting no error, got", err)
	}
	return string(bytes)
}

func export(testfile string, loopCount int, delayMultiplier float64, width int) viz.Image {
	img := viz.Image{
		Filename:        testData + testfile,
		ExportFilename:  "/tmp/img_test.sh",
		LoopCount:       loopCount,
		DelayMultiplier: delayMultiplier,
		UserWidth:       width,
	}

	os.Args = []string{"img", "-o", img.ExportFilename,
		"-l", fmt.Sprintf("%v", img.LoopCount),
		"-s", fmt.Sprintf("%v", img.DelayMultiplier),
		"-w", fmt.Sprintf("%v", img.UserWidth),
		testData + testfile,
	}
	main() //invoke main to test flag parsing as well.

	return img
}

func validate(expected string, got viz.Image, t *testing.T) {
	if read(testData+expected, t) != read(got.ExportFilename, t) {
		t.Fatalf("expected: %v, got: %v; params: loopCount=%v, delayMultiplier=%v, userWidth=%v",
			expected, got.ExportFilename, got.LoopCount, got.DelayMultiplier, got.UserWidth)
	}
}

func TestStaticImage(t *testing.T) {
	img := export("color_matrix.png", 1, 1.0, 80)
	validate("color_matrix.sh", img, t)
}

func TestGIF(t *testing.T) {
	// Override Size() because the Unix system calls in
	// terminal.GetSize() fail with "operation not permitted"
	// when executed in the test environment
	terminal.Size = func() (int, int, error) {
		return math.MaxInt32, math.MaxInt32, nil
	}
	// Test different kinds of GIF disposals.
	for _, d := range []string{"Unspecified", "None", "NoneTransparency", "Background"} {
		img := export("disposal"+d+".gif", 1, 1.0, 0)
		validate("disposal"+d+".sh", img, t)
	}

	// Test all parameters
	img := export("disposalNone.gif", 3, 2, 60)
	validate("all.sh", img, t)
}

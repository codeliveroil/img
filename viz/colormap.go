// Copyright (c) 2018 codeliveroil. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package viz

import (
	"fmt"
	"image/color"
)

// colors is a palette where the index corresponds to
// a 8 bit terminal color palette (i.e. 256 colors)
var Colors color.Palette = make([]color.Color, 256)
var _ = InitColors(256)

// newColor returns an RGBA instance.
func newColor(R, G, B uint8) color.Color {
	return color.RGBA{R: R, G: G, B: B, A: 255}
}

// newColor from bitmap takes in a byte and reads the first three bits as
// R,G,B and returns a color accordingly, replacing R/G/B with componentVal
func newColorFromBitmap(setBits uint8, componentVal uint8) color.Color {
	c := make([]uint8, 3)
	for i, bit := range []uint8{1, 2, 4} {
		if setBits&bit == bit {
			c[i] = componentVal
		}
	}
	return newColor(c[0], c[1], c[2])
}

// InitColors initializes 'Colors' (the terminal emulator color
// palette) with 'num' colors. Accepted num values are
// 8, 16, 256.
// Returns a boolean indicating success or failure.
// If 'num' is not one of the accepted values, 256 colors
// are initialized and false is returned.
func InitColors(num int) bool {
	// Primary 3 bit colors
	for i := 0; i <= 6; i++ {
		Colors[i] = newColorFromBitmap(uint8(i), 128)
	}
	Colors[7] = newColor(192, 192, 192)
	if num == 8 {
		return true
	}

	// Bright version of primary colors
	Colors[8] = newColor(128, 128, 128)
	for i := 9; i <= 15; i++ {
		Colors[i] = newColorFromBitmap(uint8(i-8), 255)
	}
	if num == 16 {
		return true
	}

	// Other colors
	pattern := []uint8{0, 95, 135, 175, 215, 255}
	i := 16
	for r := 0; r < len(pattern); r++ {
		for g := 0; g < len(pattern); g++ {
			for b := 0; b < len(pattern); b++ {
				Colors[i] = newColor(pattern[r], pattern[g], pattern[b])
				i++
			}
		}
	}

	// Remaining grays
	for i := 232; i <= 255; i++ {
		c := uint8((i-232)*10 + 8)
		Colors[i] = newColor(c, c, c)
	}

	return num == 256
}

// Draw renders all 256 colors to stdout for
// debugging purposes
func Draw() {
	for i, _ := range Colors {
		fmt.Print(fmt.Sprintf("\x1b[48;5;%vm   \x1b[0m", i))
		if i%32 == 0 {
			fmt.Println("")
		}
	}
}

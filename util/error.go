// Copyright (c) 2018 codeliveroil. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package util

import (
	"fmt"
	"os"
)

// Check prints an error message and exits
// if err is not nil.
func Check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

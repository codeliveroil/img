// Copyright (c) 2018 codeliveroil. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package terminal

import systerm "golang.org/x/crypto/ssh/terminal"

// Size returns the dimensions of the terminal.
// This function can be overriden for test cases
// as system calls fail with "operation not supported"
// in test environments
var Size = func() (width int, height int, err error) {
	return systerm.GetSize(0)
}

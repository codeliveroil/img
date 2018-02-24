// Copyright (c) 2018 codeliveroil. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package terminal

import (
	"io"
	"os"
	"os/exec"

	systerm "golang.org/x/crypto/ssh/terminal"
)

// Size returns the dimensions of the terminal.
// This function can be overriden for test cases
// as system calls fail with "operation not supported"
// in test environments
var Size = func() (width int, height int, err error) {
	return systerm.GetSize(0)
}

// LineUp moves the cursor one line up.
func LineUp() error {
	return runCommand(nil, "tput", "cuu1")
}

// runCommand runs the command and the stdout is redirected
// to the given stdout (or OS's stdout if given is nil) and
// the OS's stderr.
func runCommand(stdout io.Writer, command string, args ...string) (exitCode error) {
	cmd := exec.Command(command, args...)
	if stdout == nil {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = stdout
	}
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

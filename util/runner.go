// Copyright (c) 2018 codeliveroil. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package util

import (
	"io"
	"os"
	"os/exec"
)

// RunCommand runs the command and the stdout is redirected
// to the given stdout (or OS's stdout if given is nil) and
// the OS's stderr.
func RunCommand(stdout io.Writer, command string, args ...string) (exitCode error) {
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

// StdWriter is a simple writer to record
// output from stdout/stderr.
type StdWriter struct {
	Output []string
}

func (s *StdWriter) Write(p []byte) (n int, err error) {
	l := len(p)
	if l == 0 {
		return 0, nil
	}
	sub := 0
	if p[l-1] == '\n' {
		sub = 1
	}
	s.Output = append(s.Output, string(p[:len(p)-sub]))
	return len(p), nil
}

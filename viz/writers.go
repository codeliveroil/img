// Copyright (c) 2018 codeliveroil. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package viz

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/codeliveroil/img/util"
)

// Writer can be used to forward the output of the image
// conversion to a given destination (e.g. stdout vs file).
type Writer interface {
	//Write writes the string to the destination.
	Write(str string) error
	//Moves the cursor up one line 'count' times.
	LineUp(count int) error
	//Sleep sleeps for the specified time.
	Sleep(delayMS int) error
	//Close closes the writer.
	Close() error
}

// NewFileWriter returns a file writer.
func NewFileWriter(filename string) (*FileWriter, error) {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		if err := os.Remove(filename); err != nil {
			return nil, err
		}
	}

	f, _ := os.Open(filename)
	flags := os.O_CREATE | os.O_WRONLY
	f, err := os.OpenFile(filename, flags, 0755)
	if err != nil {
		return nil, err
	}
	w := bufio.NewWriter(f)
	fw := &FileWriter{file: f, writer: w}
	return fw, fw.Write("echo -n '")
}

// File writer writes strings to a file enclosed
// in one giant echo.
type FileWriter struct {
	file       *os.File
	writer     *bufio.Writer
	writeError error
}

func (w *FileWriter) Write(str string) error {
	if w.writeError == nil {
		_, w.writeError = w.writer.WriteString(str)
	}
	return w.writeError
}

func (w *FileWriter) LineUp(count int) error {
	w.Write("'\n") //close echo
	for i := 0; i < count; i++ {
		w.Write("tput cuu1\n")
	}
	w.Write("echo -n '")
	return w.writeError
}

func (w *FileWriter) Sleep(delayMS int) error {
	w.Write("'\n") //close echo
	w.Write(fmt.Sprintf("sleep %v\n", float64(delayMS)/1000))
	w.Write("echo -n '")
	return w.writeError
}

func (w *FileWriter) Close() error {
	if err := w.Write("'"); err != nil {
		return err
	}

	if err := w.writer.Flush(); err != nil {
		return err
	}
	if err := w.file.Close(); err != nil {
		return err
	}
	return nil
}

// StdWriter writes strings to stdout.
type StdWriter struct {
	b bytes.Buffer
}

func (w *StdWriter) Write(str string) error {
	w.b.WriteString(str)
	return nil
}

func (w *StdWriter) LineUp(count int) error {
	fmt.Printf(w.b.String())
	w.b.Reset()

	for i := 0; i < count; i++ {
		err := util.RunCommand(nil, "tput", "cuu1")
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *StdWriter) Sleep(delayMS int) error {
	time.Sleep(time.Millisecond * time.Duration(delayMS))
	return nil
}

func (w *StdWriter) Close() error {
	fmt.Printf(w.b.String())
	return nil
}

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

	"github.com/codeliveroil/img/terminal"
)

// Canvas is the destination (e.g. stdout vs file) where the
// image will be rendered.
type Canvas interface {
	// Paint renders two pixels at a time - the top (y) and the
	// bottom (y+1) ones.
	Paint(topColor, bottomColor uint8) error
	// NewLine moves the cursor to the next line
	NewLine() error
	// Moves the cursor up one line 'count' times.
	LineUp(count int) error
	// Sleep sleeps for the specified time.
	Sleep(delayMS int) error
	// Close closes the canvas.
	Close() error
}

// NewFileCanvas returns a FileCanvas.
func NewFileCanvas(filename string) (*FileCanvas, error) {
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
	fc := &FileCanvas{file: f, writer: w}
	return fc, fc.write("echo -n '")
}

// FileCanvas exports the image to a shell script
// that can be executed when required, to render the image.
type FileCanvas struct {
	file       *os.File
	writer     *bufio.Writer
	writeError error
}

func (fc *FileCanvas) write(str string) error {
	if fc.writeError == nil {
		_, fc.writeError = fc.writer.WriteString(str)
	}
	return fc.writeError
}

func (fc *FileCanvas) Paint(topColor, bottomColor uint8) error {
	return fc.write(makeTwoPixels(topColor, bottomColor))
}

func (fc *FileCanvas) NewLine() error {
	return fc.write("\n")
}

func (fc *FileCanvas) LineUp(count int) error {
	fc.write("'\n") //close echo
	for i := 0; i < count; i++ {
		fc.write("tput cuu1\n")
	}
	fc.write("echo -n '")
	return fc.writeError
}

func (fc *FileCanvas) Sleep(delayMS int) error {
	fc.write("'\n") //close echo
	fc.write(fmt.Sprintf("sleep %v\n", float64(delayMS)/1000))
	fc.write("echo -n '")
	return fc.writeError
}

func (fc *FileCanvas) Close() error {
	if err := fc.write("'"); err != nil {
		return err
	}

	if err := fc.writer.Flush(); err != nil {
		return err
	}
	if err := fc.file.Close(); err != nil {
		return err
	}
	return nil
}

// StdoutCanvas renders the image to stdout.
type StdoutCanvas struct {
	b bytes.Buffer
}

func (sc *StdoutCanvas) Paint(topColor, bottomColor uint8) error {
	sc.b.WriteString(makeTwoPixels(topColor, bottomColor))
	return nil
}

func (sc *StdoutCanvas) NewLine() error {
	sc.b.WriteString("\n")
	return nil
}

func (sc *StdoutCanvas) LineUp(count int) error {
	fmt.Printf(sc.b.String())
	sc.b.Reset()

	for i := 0; i < count; i++ {
		err := terminal.LineUp()
		if err != nil {
			return err
		}
	}
	return nil
}

func (sc *StdoutCanvas) Sleep(delayMS int) error {
	time.Sleep(time.Millisecond * time.Duration(delayMS))
	return nil
}

func (sc *StdoutCanvas) Close() error {
	fmt.Printf(sc.b.String())
	return nil
}

func makeTwoPixels(topColor, bottomColor uint8) string {
	//Use the 'lower half block' character (▄) for drawing as opposed to the 'upper half block' (▀) because if the
	//terminal character height is odd then the terminal aligns ▀ one line below the top rendering a background shade
	//on the top line of the character
	return fmt.Sprintf("\x1b[48;5;%vm\x1b[38;5;%vm▄\x1b[0m", topColor, bottomColor)
}

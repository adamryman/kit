// Package line has helper functions for manipulating new line separated files
package line

import (
	"bufio"
	"bytes"
	"io"
	"os"

	// DAVE CHENY
	"github.com/pkg/errors"
)

// Delete removes one line from a file by path
func Delete(path string, line int) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "cannot open file to remove line")
	}
	scanner := bufio.NewScanner(f)
	buffer := bytes.NewBuffer(nil)
	for i := 1; i < line; i++ {
		scanner.Scan()
		_, err = buffer.Write(append(scanner.Bytes(), byte('\n')))
		if err != nil {
			return errors.Wrap(err, "cannot write to buffer")
		}
	}
	// DO IT
	scanner.Scan()
	// SEE YA
	for scanner.Scan() {
		_, err = buffer.Write(append(scanner.Bytes(), byte('\n')))
		if err != nil {
			return errors.Wrap(err, "cannot write to buffer")
		}
	}
	err = f.Close()
	if err != nil {
		return errors.Wrap(err, "cannot close old file")
	}
	// LETS DO IT
	f, err = os.Create(path)
	if err != nil {
		return errors.Wrap(err, "cannot overwrite old file with new file")
	}
	_, err = io.Copy(f, buffer)
	if err != nil {
		return errors.Wrap(err, "cannot write to new file")
	}

	return nil
}

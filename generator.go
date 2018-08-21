// Copyright (c) 2018 Beta Kuang
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package csv

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
)

// NewGenerator creates and returns a new generator with the given settings.
func NewGenerator(settings ...Setting) *Generator {
	var g = &Generator{
		rule: defaultRule,
	}
	for _, setting := range settings {
		setting(&g.rule)
	}

	g.buf = bytes.NewBuffer(nil)
	g.w = bufio.NewWriter(g.rule.encoding.NewEncoder().Writer(g.buf))
	return g
}

// A Generator generates a new CSV document.
type Generator struct {
	rule rule
	buf  *bytes.Buffer
	w    *bufio.Writer

	finished bool
}

// Write writes a record row to the end of the document.
//
// If Finish has been called, Write returns an error.
func (g *Generator) Write(record []string) error {
	if g.finished {
		return fmt.Errorf("csv: Generator has been finished")
	}

	var err = g.writeRecord(record)
	if err != nil {
		return g.error(err)
	}
	return nil
}

// WriteAll writes all the rows in records to the end of the document.
//
// If Finish has been called, WriteAll returns an error.
func (g *Generator) WriteAll(records [][]string) error {
	if g.finished {
		return fmt.Errorf("csv: Generator has been finished")
	}

	var err error
	for _, record := range records {
		err = g.writeRecord(record)
		if err != nil {
			return g.error(err)
		}
	}
	return nil
}

func (g *Generator) error(err error) error {
	return fmt.Errorf("csv: Generator failed: %s", err.Error())
}

func (g *Generator) writeRecord(record []string) error {
	if g.w.Buffered() > 0 {
		// Write a line end if the buffer is not empty.
		_, err := g.w.WriteRune('\n')
		if err != nil {
			return err
		}
	}

	var err error
	for i := 0; i < len(record); i++ {
		var field = record[i]
		err = g.writeField(field)
		if err != nil {
			return err
		}

		if i < len(record)-1 {
			err = g.writeSeparator()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Generator) writeField(field string) error {
	// Prefix.
	if g.rule.prefix != noRune {
		_, err := g.w.WriteRune(g.rule.prefix)
		if err != nil {
			return err
		}
	}

	if strings.ContainsAny(field, "\"\n") || strings.ContainsRune(field, g.rule.separator) {
		var escaped = fmt.Sprintf(`"%s"`, strings.Replace(field, "\"", "\"\"", -1))
		_, err := g.w.WriteString(escaped)
		if err != nil {
			return err
		}
	} else {
		_, err := g.w.WriteString(field)
		if err != nil {
			return err
		}
	}

	// Suffix.
	if g.rule.suffix != noRune {
		_, err := g.w.WriteRune(g.rule.suffix)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) writeSeparator() error {
	_, err := g.w.WriteRune(g.rule.separator)
	return err
}

// Finish finishes writing to the generator and returns data of the document.
//
// After calling Finish, the generator can no longer be written. Any call to
// Write and WriteAll will return an error.
func (g *Generator) Finish() ([]byte, error) {
	g.finished = true

	var err = g.w.Flush()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(g.buf)
	if err != nil {
		return nil, err
	}
	return data, nil
}

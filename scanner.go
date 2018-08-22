// Copyright (c) 2018 Beta Kuang
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package csv

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/transform"
)

// NewScanner creates and returns a new scanner from a byte slice with the given settings.
func NewScanner(data []byte, settings ...Setting) (*Scanner, error) {
	var s = &Scanner{
		rule: defaultRule,
	}
	for _, setting := range settings {
		setting(&s.rule)
	}

	s.f = bufio.NewReader(transform.NewReader(bytes.NewReader(data), s.rule.encoding.NewDecoder()))
	var err = s.next()
	if err != nil {
		return nil, err
	}
	return s, nil
}

// A Scanner scans a CSV document and returns the scanned header and rows.
type Scanner struct {
	f    *bufio.Reader
	rule rule

	line     string
	lineNo   int
	pos      int
	c        rune
	eof      bool
	lastLine bool
}

// Setting applies settings for s.
func (s *Scanner) Setting(settings ...Setting) {
	for _, setting := range settings {
		setting(&s.rule)
	}
}

// Scan scans the next row from the CSV document.
//
// If an error occurs, row will be returned as nil.
//
// If there is no more row to be scanned, io.EOF will be returned.
func (s *Scanner) Scan() (row []string, err error) {
	if s.eof {
		return nil, io.EOF
	}

	row, err = s.scanRecord()
	if err != nil {
		return nil, s.error(err)
	}
	if s.eof {
		err = io.EOF
	}
	return
}

// ScanAll scans the rest rows of the CSV document.
//
// If an error occurs, rows will be returned as nil.
func (s *Scanner) ScanAll() (rows [][]string, err error) {
	rows = make([][]string, 0)
	for !s.eof {
		row, err := s.scanRecord()
		if err != nil {
			return nil, s.error(err)
		}
		rows = append(rows, row)
	}
	return
}

// error wraps a scanning error with the current position in the CSV document.
func (s *Scanner) error(err error) error {
	return fmt.Errorf("csv: Scanner failed at line %d, pos %d: %v", s.lineNo, s.pos, err)
}

// next moves to the next rune in the document.
func (s *Scanner) next() error {
	if s.pos >= len([]rune(s.line))-1 {
		return s.nextLine()
	}
	s.pos++
	s.c = []rune(s.line)[s.pos]
	return nil
}

// nextLine reads the next line into s.line, and updates s.c and s.pos to the
// first rune of the new line.
//
// If the new line is the last line of the document, s.lastLine will be set
// true. If the last line is empty, s.eof will be set true.
func (s *Scanner) nextLine() error {
	s.lineNo++
	var err = s.readNextLine()
	if err != nil {
		return err
	}

	// Omit comments and empty lines if necessary.
	for !s.lastLine && s.shouldOmitLine(s.line) {
		s.lineNo++
		err = s.readNextLine()
		if err != nil {
			return err
		}
	}

	s.pos = 0
	if len([]rune(s.line)) <= 0 && s.lastLine {
		s.eof = true
		s.c = noRune
		if !s.rule.allowEndingLineBreakInLastRecord {
			return fmt.Errorf("unexpected empty line")
		}
		return nil
	}
	s.c = []rune(s.line)[0]
	return nil
}

func (s *Scanner) shouldOmitLine(line string) bool {
	// Empty line (only with a line break).
	if line == "\n" && s.rule.omitEmptyLine {
		return true
	}
	// Comment.
	if s.rule.comment != noRune && len([]rune(line)) > 0 && []rune(line)[0] == s.rule.comment {
		return true
	}
	return false
}

func (s *Scanner) readNextLine() error {
	var err error
	s.line, err = s.f.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			s.lastLine = true
		} else {
			return err
		}
	}
	return nil
}

func (s *Scanner) scanRecord() ([]string, error) {
	var fields = make([]string, 0)
	field, err := s.scanField()
	if err != nil {
		return nil, err
	}
	fields = append(fields, field)

	for !s.eof && !s.isLineEnd(s.c) {
		_, err := s.scanCOMMA()
		if err != nil {
			return nil, err
		}

		field, err := s.scanField()
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}

	err = s.nextLine()
	if err != nil {
		return nil, err
	}
	return fields, nil
}

// scanField scans and returns a field.
//
// If the field starts with a quote, scanField scans until a matching quote is
// found. If the field does not start with a quote, scanField scans until a
// separator or line end is found.
//
// If the end of a field could not be found, an error will be returned.
func (s *Scanner) scanField() (string, error) {
	if s.rule.omitLeadingSpace {
		_, err := s.scanSPACE()
		if err != nil {
			return "", err
		}
	}

	if s.rule.prefix != noRune {
		if s.c == s.rule.prefix {
			var err = s.next()
			if err != nil {
				return "", err
			}
		} else {
			return "", fmt.Errorf("prefix not found")
		}
	}

	var field string
	var err error
	if s.isQuote(s.c) {
		field, err = s.scanEscaped()
		if err != nil {
			return "", err
		}
		if s.rule.suffix != noRune {
			if s.c != s.rule.suffix {
				return "", fmt.Errorf("suffix not found")
			}
			err = s.next()
			if err != nil {
				return "", err
			}
		}
	} else {
		field, err = s.scanNonEscaped()
		if err != nil {
			return "", err
		}
	}

	if s.rule.omitTrailingSpace {
		field = strings.TrimRightFunc(field, s.isSpace)
		_, err := s.scanSPACE()
		if err != nil {
			return "", err
		}
	}

	return field, nil
}

func (s *Scanner) scanEscaped() (string, error) {
	leadingQuote, err := s.scanQUOTE()
	if err != nil {
		return "", err
	}

	var escaped string
	var foundFirstQuote = false
	for !s.eof {
		if s.isQuote(s.c) {
			if string(s.c) != leadingQuote {
				if foundFirstQuote {
					return escaped, nil
				}
				escaped += string(s.c)
				err = s.next()
				if err != nil {
					return "", err
				}
				continue
			}

			// s.c == leading quote, escape or field end.
			if !foundFirstQuote {
				foundFirstQuote = true
				err = s.next()
				if err != nil {
					return "", err
				}
			} else {
				foundFirstQuote = false
				escaped += string(s.c)
				err = s.next()
				if err != nil {
					return "", err
				}
			}
		} else {
			if foundFirstQuote {
				return escaped, nil
			}
			escaped += string(s.c)
			var err = s.next()
			if err != nil {
				return "", err
			}
		}
	}

	if foundFirstQuote {
		return escaped, nil
	}
	return "", fmt.Errorf("trailing quote not found")
}

func (s *Scanner) scanNonEscaped() (string, error) {
	if (s.isComma(s.c) || s.isLineEnd(s.c) || s.eof) && !s.rule.allowEmptyField {
		return "", fmt.Errorf("unexpected empty field, expect text")
	}
	if s.isQuote(s.c) {
		return "", fmt.Errorf("unexpected character '%s', expect text", string(s.c))
	}

	var nonEscaped string
	for !s.eof && !s.isLineEnd(s.c) && !s.isComma(s.c) && (s.rule.suffix == noRune || s.c != s.rule.suffix) {
		nonEscaped += string(s.c)
		var err = s.next()
		if err != nil {
			return "", err
		}
	}

	// Suffix.
	if s.rule.suffix != noRune {
		if s.c != s.rule.suffix {
			return "", fmt.Errorf("suffix not found")
		}
		var err = s.next()
		if err != nil {
			return "", err
		}
	}
	return nonEscaped, nil
}

// scanCOMMA scans a separator. A separator is a comma, or other rune as set
// with the Separator() setting.
//
// If no separator is found, an error will be returned.
func (s *Scanner) scanCOMMA() (string, error) {
	if s.c != s.rule.separator {
		return "", fmt.Errorf("unexpected character '%s', expect %s", string(s.c), string(s.rule.separator))
	}
	var comma = string(s.c)
	var err = s.next()
	if err != nil {
		return "", err
	}
	return comma, nil
}

// scanCRLF scans and returns a line end.
func (s *Scanner) scanCRLF() (string, error) {
	if !s.isLineEnd(s.c) {
		return "", fmt.Errorf("unexpected character '%s', expect line end", string(s.c))
	}
	var lineEnd = string(s.c)
	var err = s.next()
	if err != nil {
		return "", err
	}
	return lineEnd, nil
}

// scanQUOTE scans and returns a quote. By default, both double and single quote
// are allowed. This can be changed with the AllowSingleQuote() setting.
func (s *Scanner) scanQUOTE() (string, error) {
	if !s.isQuote(s.c) {
		return "", fmt.Errorf("unexpected character '%s', expect quote", string(s.c))
	}
	var quote = string(s.c)
	var err = s.next()
	if err != nil {
		return "", err
	}
	return quote, nil
}

// scanSPACE scans while the current rune is a space.
func (s *Scanner) scanSPACE() (string, error) {
	var spaces string
	for !s.eof && s.isSpace(s.c) {
		spaces += string(s.c)
		var err = s.next()
		if err != nil {
			return "", err
		}
	}
	return spaces, nil
}

func (s *Scanner) isQuote(c rune) bool {
	if c == '"' {
		return true
	}
	if s.rule.allowSingleQuote && c == '\'' {
		return true
	}
	return false
}

func (s *Scanner) isLineEnd(c rune) bool {
	return c == '\n'
}

func (s *Scanner) isComma(c rune) bool {
	return c == s.rule.separator
}

func (s *Scanner) isSpace(c rune) bool {
	switch c {
	case '\t', '\v', '\f', ' ', 0x85, 0xA0:
		return true
	}
	return false
}

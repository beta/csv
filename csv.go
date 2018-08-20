// Copyright (c) 2018 Beta Kuang
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package csv implements a parser and generator for CSV-like documents.
package csv

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

const noRune = '\x00'

type rule struct {
	// Common rules.
	encoding  encoding.Encoding
	separator rune

	// Scanner rules.
	allowSingleQuote  bool
	allowEmptyField   bool
	omitLeadingSpace  bool
	omitTrailingSpace bool
	comment           rune
	header            bool

	// Unmarshaler rules.
	validators map[string]func(interface{}) bool
}

var defaultRule = rule{
	// Common rules.
	encoding:  unicode.UTF8,
	separator: ',',

	// Scanner rules.
	allowSingleQuote:  true,
	allowEmptyField:   true,
	omitLeadingSpace:  true,
	omitTrailingSpace: true,
	comment:           noRune,
	header:            false,

	// Unmarshaler rules.
	validators: nil,
}

// A Setting provides information on how documents should be parsed.
type Setting func(*rule)

//==============================================================================
// Common settings.
//==============================================================================

// Encoding sets the character encoding used while reading and writing a document.
func Encoding(enc encoding.Encoding) Setting {
	return func(r *rule) {
		r.encoding = enc
	}
}

// Separator sets the separator used to separate fields while reading and writing a document.
func Separator(sep rune) Setting {
	return func(r *rule) {
		r.separator = sep
	}
}

// RFC4180 sets the parser and generator to work in the exact way as
// described in RFC 4180.
func RFC4180() Setting {
	return func(r *rule) {
		r.allowSingleQuote = false
		r.allowEmptyField = false
		r.omitLeadingSpace = false
		r.omitTrailingSpace = false
		r.comment = '\x00'
		r.separator = ','
	}
}

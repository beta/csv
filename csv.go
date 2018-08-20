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
	prefix    rune
	suffix    rune

	// Scanner rules.
	allowSingleQuote  bool
	allowEmptyField   bool
	omitLeadingSpace  bool
	omitTrailingSpace bool
	comment           rune

	// Unmarshaler rules.
	validators map[string]func(interface{}) bool

	// Marshaler rules.
	writeHeader bool
}

var defaultRule = rule{
	// Common rules.
	encoding:  unicode.UTF8,
	separator: ',',
	prefix:    noRune,
	suffix:    noRune,

	// Scanner rules.
	allowSingleQuote:  true,
	allowEmptyField:   true,
	omitLeadingSpace:  true,
	omitTrailingSpace: true,
	comment:           noRune,

	// Unmarshaler rules.
	validators: nil,

	// Marshaler rules.
	writeHeader: true,
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

// Prefix sets the prefix of every field while reading and writing a document.
func Prefix(prefix rune) Setting {
	return func(r *rule) {
		r.prefix = prefix
	}
}

// Suffix sets the suffix of every field when reading and writing a document.
func Suffix(suffix rune) Setting {
	return func(r *rule) {
		r.suffix = suffix
	}
}

// RFC4180 sets the parser and generator to work in the exact way as
// described in RFC 4180.
func RFC4180() Setting {
	return func(r *rule) {
		// Common rules.
		r.separator = ','
		r.prefix = noRune
		r.suffix = noRune

		// Scanner rules.
		r.allowSingleQuote = false
		r.allowEmptyField = false
		r.omitLeadingSpace = false
		r.omitTrailingSpace = false
		r.comment = noRune
	}
}

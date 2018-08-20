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
	encoding     encoding.Encoding
	separator    rune
	headerPrefix rune
	headerSuffix rune
	fieldPrefix  rune
	fieldSuffix  rune

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
	encoding:     unicode.UTF8,
	separator:    ',',
	headerPrefix: noRune,
	headerSuffix: noRune,
	fieldPrefix:  noRune,
	fieldSuffix:  noRune,

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

// HeaderPrefix sets the prefix of every header name while reading and writing a document.
//
// This setting will also set Header(true).
func HeaderPrefix(prefix rune) Setting {
	return func(r *rule) {
		r.header = true
		r.headerPrefix = prefix
	}
}

// HeaderSuffix sets the suffix of every header name while reading and writing to a document.
//
// This setting will also set Header(true).
func HeaderSuffix(suffix rune) Setting {
	return func(r *rule) {
		r.header = true
		r.headerSuffix = suffix
	}
}

// FieldPrefix sets the prefix of every field while reading and writing a document.
func FieldPrefix(prefix rune) Setting {
	return func(r *rule) {
		r.fieldPrefix = prefix
	}
}

// FieldSuffix sets the suffix of every field when reading and writing a document.
func FieldSuffix(suffix rune) Setting {
	return func(r *rule) {
		r.fieldSuffix = suffix
	}
}

// RFC4180 sets the parser and generator to work in the exact way as
// described in RFC 4180.
func RFC4180() Setting {
	return func(r *rule) {
		// Common rules.
		r.separator = ','
		r.headerPrefix = noRune
		r.headerSuffix = noRune
		r.fieldPrefix = noRune
		r.fieldSuffix = noRune

		// Scanner rules.
		r.allowSingleQuote = false
		r.allowEmptyField = false
		r.omitLeadingSpace = false
		r.omitTrailingSpace = false
		r.comment = noRune
	}
}

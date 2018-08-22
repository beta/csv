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
	allowSingleQuote                 bool
	allowEmptyField                  bool
	allowEndingLineBreakInLastRecord bool
	omitLeadingSpace                 bool
	omitTrailingSpace                bool
	omitEmptyLine                    bool
	comment                          rune

	// Unmarshaler and marshaler common rules.
	headerPrefix rune
	headerSuffix rune
	fieldPrefix  rune
	fieldSuffix  rune

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
	allowSingleQuote:                 true,
	allowEmptyField:                  true,
	allowEndingLineBreakInLastRecord: true,
	omitLeadingSpace:                 true,
	omitTrailingSpace:                true,
	omitEmptyLine:                    true,
	comment:                          noRune,

	// Unmarshaler and marshaler common rules.
	headerPrefix: noRune,
	headerSuffix: noRune,
	fieldPrefix:  noRune,
	fieldSuffix:  noRune,

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

//==============================================================================
// Scanner settings.
//==============================================================================

// AllowSingleQuote sets whether single quotes are allowed while reading a document.
func AllowSingleQuote(v bool) Setting {
	return func(r *rule) {
		r.allowSingleQuote = v
	}
}

// AllowEmptyField sets whether empty fields are allowed while reading a document.
func AllowEmptyField(v bool) Setting {
	return func(r *rule) {
		r.allowEmptyField = v
	}
}

// AllowEndingLineBreakInLastRecord sets whether the last record may have an
// ending line break while reading a document.
func AllowEndingLineBreakInLastRecord(v bool) Setting {
	return func(r *rule) {
		r.allowEndingLineBreakInLastRecord = v
	}
}

// OmitLeadingSpace sets whether the leading spaces of fields should be omitted
// while reading a document.
func OmitLeadingSpace(v bool) Setting {
	return func(r *rule) {
		r.omitLeadingSpace = v
	}
}

// OmitTrailingSpace sets whether the trailing spaces of fields should be
// omitted while reading a document.
func OmitTrailingSpace(v bool) Setting {
	return func(r *rule) {
		r.omitTrailingSpace = v
	}
}

// OmitEmptyLine sets whether empty lines should be omitted while reading a document.
func OmitEmptyLine(v bool) Setting {
	return func(r *rule) {
		r.omitEmptyLine = v
	}
}

// Comment sets the leading rune of comments used while reading a document.
func Comment(comment rune) Setting {
	return func(r *rule) {
		r.comment = comment
	}
}

//==============================================================================
// Unmarshaler and marshaler common settings.
//==============================================================================

// HeaderPrefix sets the prefix rune of header names while unmarshaling and
// marshaling a document.
//
// If a header prefix is set, the Prefix setting will be ignored while reading
// and writing the header row, but will still be used for fields.
func HeaderPrefix(prefix rune) Setting {
	return func(r *rule) {
		r.headerPrefix = prefix
	}
}

// HeaderSuffix sets the suffix rune of header names while unmarshaling and
// marshaling a document.
//
// If a header suffix is set, the Suffix setting will be ignored while reading
// and writing the header row, but will still be used for fields.
func HeaderSuffix(suffix rune) Setting {
	return func(r *rule) {
		r.headerSuffix = suffix
	}
}

// FieldPrefix sets the prefix rune of fields while unmarshaling and marshaling
// a document.
//
// If a field prefix is set, the Prefix setting will be ignored while reading
// and writing fields, but will still be used for the header.
func FieldPrefix(prefix rune) Setting {
	return func(r *rule) {
		r.fieldPrefix = prefix
	}
}

// FieldSuffix sets the suffix rune of fields while unmarshaling and marshaling
// a document.
//
// If a field suffix is set, the Suffix setting will be ignored while reading
// and writing fields, but will still be used for the header.
func FieldSuffix(suffix rune) Setting {
	return func(r *rule) {
		r.fieldSuffix = suffix
	}
}

//==============================================================================
// Unmarshaler settings.
//==============================================================================

// Validator adds a new validator functions for validating a CSV value while
// unmarshaling a document.
func Validator(name string, validator func(interface{}) bool) Setting {
	return func(r *rule) {
		if r.validators == nil {
			r.validators = make(map[string]func(interface{}) bool)
		}
		r.validators[name] = validator
	}
}

//==============================================================================
// Marshaler settings.
//==============================================================================

// WriteHeader sets whether to output the header row while writing the document.
func WriteHeader(v bool) Setting {
	return func(r *rule) {
		r.writeHeader = v
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
		r.allowEndingLineBreakInLastRecord = true
		r.omitLeadingSpace = false
		r.omitTrailingSpace = false
		r.omitEmptyLine = false
		r.comment = noRune

		// Unmarshaler and marshaler common settings.
		r.headerPrefix = noRune
		r.headerSuffix = noRune
		r.fieldPrefix = noRune
		r.fieldSuffix = noRune
	}
}

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

type rule struct {
	// Common rules.
	encoding  encoding.Encoding
	separator rune

	// Scanner rules.
	allowSingleQuote  bool
	allowEmptyField   bool
	omitLeadingSpace  bool
	omitTrailingSpace bool
	allowComment      bool
	comment           rune
	header            bool

	// Unmarshaler rules.
	validators map[string]func(interface{}) bool
}

var defaultRule = rule{
	encoding:  unicode.UTF8,
	separator: ',',

	allowSingleQuote:  true,
	allowEmptyField:   true,
	omitLeadingSpace:  true,
	omitTrailingSpace: true,
	allowComment:      true,
	comment:           ';',
	header:            false,

	validators: nil,
}

// A Setting provides information on how documents should be parsed.
type Setting func(*rule)

var (
	// Encoding sets the character encoding used while reading and writing a document.
	Encoding = func(enc encoding.Encoding) Setting {
		return func(r *rule) {
			r.encoding = enc
		}
	}

	// AllowSingleQuote sets whether single quotes are allowed while reading a document.
	AllowSingleQuote = func(v bool) Setting {
		return func(r *rule) {
			r.allowSingleQuote = v
		}
	}
	// AllowEmptyField sets whether empty fields are allowed while reading a document.
	AllowEmptyField = func(v bool) Setting {
		return func(r *rule) {
			r.allowEmptyField = v
		}
	}
	// OmitLeadingSpace sets whether the leading spaces of fields should be omitted while reading a document.
	OmitLeadingSpace = func(v bool) Setting {
		return func(r *rule) {
			r.omitLeadingSpace = v
		}
	}
	// OmitTrailingSpace sets whether the trailing spaces of fields should be omitted while reading a document.
	OmitTrailingSpace = func(v bool) Setting {
		return func(r *rule) {
			r.omitTrailingSpace = v
		}
	}
	// AllowComment sets whether comments are allowed (and ignored) while reading a document.
	AllowComment = func(v bool) Setting {
		return func(r *rule) {
			r.allowComment = v
		}
	}
	// Comment sets the leading rune of comments used while reading a document.
	Comment = func(comment rune) Setting {
		return func(r *rule) {
			r.comment = comment
		}
	}
	// Header sets whether there is a header to be read while reading the document.
	Header = func(v bool) Setting {
		return func(r *rule) {
			r.header = v
		}
	}
	// Separator sets the separator used to separate fields while reading and writing a document.
	Separator = func(sep rune) Setting {
		return func(r *rule) {
			r.separator = sep
		}
	}

	// RFC4180 sets the parser and generator to work in the exact way as
	// described in RFC 4180.
	RFC4180 = func() Setting {
		return func(r *rule) {
			r.allowSingleQuote = false
			r.allowEmptyField = false
			r.omitLeadingSpace = false
			r.omitTrailingSpace = false
			r.allowComment = false
			r.comment = '\x00'
			r.separator = ','
		}
	}
)

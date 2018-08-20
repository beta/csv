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
	encoding          encoding.Encoding
	allowSingleQuote  bool
	allowEmptyField   bool
	omitLeadingSpace  bool
	omitTrailingSpace bool
	allowComment      bool
	comment           rune
	separator         rune
	header            bool
}

var defaultRule = rule{
	encoding:          unicode.UTF8,
	allowSingleQuote:  true,
	allowEmptyField:   true,
	omitLeadingSpace:  true,
	omitTrailingSpace: true,
	allowComment:      true,
	comment:           ';',
	separator:         ',',
	header:            false,
}

// A Setting provides information on how documents should be parsed.
type Setting func(*rule)

var (
	// Encoding sets the character encoding used to read or write to the document.
	Encoding = func(enc encoding.Encoding) Setting {
		return func(r *rule) {
			r.encoding = enc
		}
	}
	// AllowSingleQuote sets whether single quotes are allowed in the document.
	AllowSingleQuote = func(v bool) Setting {
		return func(r *rule) {
			r.allowSingleQuote = v
		}
	}
	// AllowEmptyField sets whether empty fields are allowed.
	AllowEmptyField = func(v bool) Setting {
		return func(r *rule) {
			r.allowEmptyField = v
		}
	}
	// OmitLeadingSpace sets whether the leading spaces of a field should be omitted.
	OmitLeadingSpace = func(v bool) Setting {
		return func(r *rule) {
			r.omitLeadingSpace = v
		}
	}
	// OmitTrailingSpace sets whether the trailing spaces of a field should be omitted.
	OmitTrailingSpace = func(v bool) Setting {
		return func(r *rule) {
			r.omitTrailingSpace = v
		}
	}
	// AllowComment sets whether comments are allowed.
	AllowComment = func(v bool) Setting {
		return func(r *rule) {
			r.allowComment = v
		}
	}
	// Comment sets the leading rune of comments.
	Comment = func(comment rune) Setting {
		return func(r *rule) {
			r.comment = comment
		}
	}
	// Separator sets the separator used to separate fields.
	Separator = func(sep rune) Setting {
		return func(r *rule) {
			r.separator = sep
		}
	}
	// Header sets whether there is a header row in the document.
	Header = func(v bool) Setting {
		return func(r *rule) {
			r.header = v
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

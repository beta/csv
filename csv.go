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
	encoding               encoding.Encoding
	allowSingleQuote       bool
	omitLeadingSpace       bool
	allowVariableRowLength bool
	comment                rune
	separator              rune
	header                 bool
}

var defaultRule = rule{
	encoding:               unicode.UTF8,
	allowSingleQuote:       true,
	omitLeadingSpace:       true,
	allowVariableRowLength: true,
	comment:                ';',
	separator:              ',',
	header:                 false,
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
	// OmitLeadingSpace sets whether the leading spaces of a field should be omitted.
	OmitLeadingSpace = func(v bool) Setting {
		return func(r *rule) {
			r.omitLeadingSpace = v
		}
	}
	// AllowVariableRowLength sets whether variable row length is allowed. If
	// set to true, rows with less field count than others will be allowed, and
	// the missing fields will be set as empty string.
	//
	// If a row has more fields than any of the rows above, empty fields will be
	// appended to all the rows to achieve the same field count.
	//
	// If the document has a header, the field count of rows must not be longer
	// than the header, or an error will be raised.
	AllowVariableRowLength = func(v bool) Setting {
		return func(r *rule) {
			r.allowVariableRowLength = v
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
			r.omitLeadingSpace = false
			r.separator = ','
			r.comment = '\x00'
		}
	}
)

// Copyright (c) 2018 Beta Kuang
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package csv

// Marshal generates a CSV document from v with the given settings.
//
// v should be an array/slice of struct or struct pointers. In these structs,
// each exported field will be marshaled as a CSV field, with the field name as
// the column header. This can be customized with a "csv" struct field tag,
// which gives the name of the field. Use "-" to omit a field from being
// marshaled. Use "-," to set the header name to "-".
//
// Below are some example of using the "csv" struct field tag.
//
//     // Field will be marshaled with "myName" as its header name.
//     Field int `csv:"myName"`
//
//     // Field is ignored.
//     Field int `csv:"-"`
//
//     // Field will be marshaled with "-" as its header name.
//     Field int `csv:"-,"`
//
// Marshal supports the following types:
//
// A boolean value will be marshaled to "true" or "false" based on its value.
//
// A floating point, integer or number value will be marshaled to the string
// representation of its value.
//
// A string value will be marshaled to the value of itself.
//
// Any type implementing encoding.TextMarshaler will be marshaled to the value
// returned by MarshalText.
//
// If Marshal encounters a field with an unsupported type, an
// UnsupportedTypeError will be returned.
//
// In order to marshal an unsupported type, a translator can be used to
// translate the value. A translator is a func(interface{}) ([]byte, error). Use
// the Translator setting to register one or more translators before marshaling.
// For example:
//
//     func TranslateIntSlice(slice interface{}) ([]byte, error) {
//         ...
//     }
//
//     csv.Marshal(..., csv.Translator("intSlice", TranslateIntSlice))
//
// To use a translator for an unsupported type, add it to the "csv" struct field
// tag. For example:
//
//     // Field will be marshaled with "myName" as its header name, and use
//     // translator with name "intSlice" to translate its value.
//     Field []int `csv:"myName,intSlice"`
//
// If a field has multiple ways to be marshaled, the order of using these ways
// is:
//
//     1. Using the translator specified in "csv" struct field tag.
//     2. Call MarshalText of the field.
//     3. Use the default way to marshal the field if it is supported.
func Marshal(v interface{}, settings ...Setting) ([]byte, error) {
	return nil, nil
}

func newMarshaler(v interface{}, settings ...Setting) *marshaler {
	return nil
}

type marshaler struct {
	//
}

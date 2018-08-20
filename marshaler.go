// Copyright (c) 2018 Beta Kuang
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package csv

//==============================================================================
// Marshaler settings.
//==============================================================================

// WriteHeader sets whether to output the header row while writing the document.
func WriteHeader(v bool) Setting {
	return func(r *rule) {
		r.writeHeader = v
	}
}

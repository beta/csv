# CSV

<a href="https://godoc.org/github.com/beta/csv">![godoc](https://img.shields.io/badge/godoc-reference-%235272B4.svg)</a>

A Golang package for reading and writing CSV-like documents.

Working in progress.

## Why another CSV package?

Golang provides an `encoding/csv` package out of the box for reading and writing standard CSV files (as described in [RFC 4180](https://tools.ietf.org/html/rfc4180)). However, not all CSV documents follow the specs. Although the built-in `csv` package provides some sort of customizability, it cannot cover all variations of these CSV formats.

Also, the built-in `csv` package does not implement a marshaler/unmarshaler (like in `json`) for generating/parsing CSV documents from/into struct instances.

This package aims to offer better customizability than the built-in `csv` package, as well as a set of easy-to-use marshaler and unmarshaler.

## Installation

`$ go get -u github.com/beta/csv`

## Getting started

TBD.

## Customization

`csv` uses setting functions for customization. The default setting used by `csv` supports a flexible variation of the standard CSV format.

- The default encoding is UTF-8.
- `,` is used as the default separator.
- Single quotes are allowed. Escaping works exactly the same as double quotes.

  `"field 1",'field 2','field 3 ''escaped''','field "4"'`

  will be parsed as

  `["field 1", "field 2", "field 3 'escaped'", "field \"4\""]`
- Empty fields are allowed.

  `field 1,,field 3`

  will be parsed as

  `["field 1", "", "field 3"]`
- An ending line break in the last record is allowed.
- Leading and trailing spaces in fields will be ignored.

  `field 1  , field 2`

  will be parsed as

  `["field 1", "field 2"]`
- Empty lines are omitted.
- The marshaler by default outputs the header row based on the `csv` tag of struct fields.

Below lists all the settings that can be used to customize the behavior of `csv`.

### Common settings

| Setting                       | Description                                                                      | Default        |
| ----------------------------- | -------------------------------------------------------------------------------- | -------------- |
| `Encoding(encoding.Encoding)` | Sets the character encoding used while reading and writing a document.           | `unicode.UTF8` |
| `Separator(rune)`             | Sets the separator used to separate fields while reading and writing a document. | `,`            |
| `Prefix(rune)`                | Sets the prefix of every field while reading and writing a document.             |                |
| `Suffix(rune)`                | Sets the suffix of every field while reading and writing a document.             |                |

### Scanner settings

| Setting                                  | Description                                                                                                                                                                                                                                                            | Default |
| ---------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| `AllowSingleQuote(bool)`                 | Sets whether single quotes are allowed while scanning a document.                                                                                                                                                                                                      | `true`  |
| `AllowEmptyField(bool)`                  | Sets whether empty fields are allowed while scanning a document.                                                                                                                                                                                                       | `true`  |
| `AllowEndingLineBreakInLastRecord(bool)` | Sets whether the last record may have an ending line break while reading a document.                                                                                                                                                                                   | `true`  |
| `OmitLeadingSpace(bool)`                 | Sets whether the leading spaces of fields should be omitted while scanning a document.                                                                                                                                                                                 | `true`  |
| `OmitTrailingSpace(bool)`                | Sets whether the trailing spaces of fields should be omitted while scanning a document.                                                                                                                                                                                | `true`  |
| `OmitEmptyLine(bool)`                    | Sets whether empty lines should be omitted while reading a document.                                                                                                                                                                                                   | `true`  |
| `Comment(rune)`                          | Sets the leading rune of comments used while scanning a document.                                                                                                                                                                                                      |         |
| `IgnoreBOM(bool)`                        | Sets whether the leading BOM (byte order mark) should be ignored while reading a document. If not, the BOM will be treated as normal content.<br>This should not be done by a csv package, but since Golang has no built-in support for BOM, a workaround is required. | `true`  |

### Unmarshaler and marshaler settings

| Setting              | Description                                                                                                                                                                                                                        | Default |
| -------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| `HeaderPrefix(rune)` | Sets the prefix rune of header names while unmarshaling and marshaling a document.<br>If a header prefix is set, the `Prefix` setting will be ignored while reading and writing the header row, but will still be used for fields. |         |
| `HeaderSuffix(rune)` | Sets the suffix rune of header names while unmarshaling and marshaling a document.<br>If a header suffix is set, the `Suffix` setting will be ignored while reading and writing the header row, but will still be used for fields. |         |
| `FieldPrefix(rune)`  | Sets the prefix rune of fields while unmarshaling and marshaling a document.<br>If a field prefix is set, the `Prefix` setting will be ignored while reading and writing fields, but will still be used for the header.            |         |
| `FieldSuffix(rune)`  | Sets the suffix rune of fields while unmarshaling and marshaling a document.<br>If a field suffix is set, the `Suffix` setting will be ignored while reading and writing fields, but will still be used for the header.            |         |

### Unmarshaler settings

| Setting                                     | Description                                                                             | Default |
| ------------------------------------------- | --------------------------------------------------------------------------------------- | ------- |
| `Validator(string, func(interface{}) bool)` | Adds a new validator function for validating a CSV value while unmarshaling a document. |         |

### Marshaler settings

| Setting             | Description                                                       | Default |
| ------------------- | ----------------------------------------------------------------- | ------- |
| `WriteHeader(bool)` | Sets whether to output the header row while writing the document. | `true`  |

All scanner settings can be used in an unmarshaler. Also, all generator settings can be used in an marshaler.

Beside the settings above, there's a special setting named `RFC4180` which applies the requirements as described in [RFC 4180](https://tools.ietf.org/html/rfc4180), including

- using `,` as the separator,
- no prefix and suffix,
- not allowing single quotes,
- not allowing empty fields,
- allowing an ending line break in the last record,
- not omitting leading and trailing spaces,
- not omitting empty lines, and
- not allowing comments.

## License

MIT
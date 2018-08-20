# CSV

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

## Settings

`csv` uses setting functions for customization. Below lists the settings supported by `csv`.

### Common settings

| Setting                       | Description                                                                      | Default        |
| ----------------------------- | -------------------------------------------------------------------------------- | -------------- |
| `Encoding(encoding.Encoding)` | Sets the character encoding used while reading and writing a document.           | `unicode.UTF8` |
| `Separator(rune)`             | Sets the separator used to separate fields while reading and writing a document. | `,`            |
| `Prefix(rune)`                | Sets the prefix of every field while reading and writing a document.             |                |
| `Suffix(rune)`                | Sets the suffix of every field while reading and writing a document.             |                |

### Scanner settings

| Setting                   | Description                                                                             | Default |
| ------------------------- | --------------------------------------------------------------------------------------- | ------- |
| `AllowSingleQuote(bool)`  | Sets whether single quotes are allowed while scanning a document.                       | `true`  |
| `AllowEmptyField(bool)`   | Sets whether empty fields are allowed while scanning a document.                        | `true`  |
| `OmitLeadingSpace(bool)`  | Sets whether the leading spaces of fields should be omitted while scanning a document.  | `true`  |
| `OmitTrailingSpace(bool)` | Sets whether the trailing spaces of fields should be omitted while scanning a document. | `true`  |
| `Comment(rune)`           | Sets the leading rune of comments used while scanning a document.                       |         |

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
- not omitting leading and trailing spaces, and
- not allowing comments.

## License

MIT
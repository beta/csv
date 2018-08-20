# CSV

A Golang package for scanning, parsing and generating CSV-like documents.

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

| Name and params               | Description                                                                            | Default        |
| ----------------------------- | -------------------------------------------------------------------------------------- | -------------- |
| `Encoding(encoding.Encoding)` | Sets the character encoding used while reading and writing a document.                 | `unicode.UTF8` |
| `AllowSingleQuote(bool)`      | Sets whether single quotes are allowed while reading a document.                       | `true`         |
| `AllowEmptyField(bool)`       | Sets whether empty fields are allowed while reading a document.                        | `true`         |
| `OmitLeadingSpace(bool)`      | Sets whether the leading spaces of fields should be omitted while reading a document.  | `true`         |
| `OmitTrailingSpace(bool)`     | Sets whether the trailing spaces of fields should be omitted while reading a document. | `true`         |
| `AllowComment(bool)`          | Sets whether comments are allowed (and ignored) while reading a document.              | `true`         |
| `Comment(rune)`               | Sets the leading rune of comments used while reading a document.                       | `;`            |
| `Header(bool)`                | Sets whether there is a header to be read while reading the document.                  | `false`        |
| `Separator(rune)`             | Sets the separator used to separate fields while reading and writing a document.       | `,`            |

To use a setting, pass it to the `New...` functions. For example:

```go
var scanner = csv.NewScanner(r, csv.AllowEmptyField(true), csv.Header(true))
```

Beside the settings above, there's a special setting named `RFC4180` which applies the requirements as described in [RFC 4180](https://tools.ietf.org/html/rfc4180), including

- not allowing single quotes,
- not allowing empty fields,
- not omitting leading and trailing spaces,
- not allowing comments, and
- using `,` as the separator.

## License

MIT
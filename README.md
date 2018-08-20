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

## License

MIT
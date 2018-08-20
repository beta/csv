// Copyright (c) 2018 Beta Kuang
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package csv_test

import (
	"strings"
	"testing"

	"github.com/beta/csv"
)

const csvStandard = `aaa,bbb,ccc
"aaa",bbb,"ccc"
aaa,"b
bb","ccc"`

const csvEscaped = `aaa,bbb,ccc
"aa""a",bbb,"ccc"
aaa,"b""
b""b","ccc"`

const csvWithHeader = `Col A,"Col B","Col
C"
aaa,bbb,ccc
"aaa",bbb,"ccc"
aaa,"b
bb","ccc"`

const csvWithCustomSeparator = `Col A|"Col B"|"Col
C"
aaa|bbb|ccc
"aaa"|bbb|"ccc"
aaa|"b
bb"|"ccc"`

const csvWithSingleQuote = `aaa,bbb,ccc
"aa'a",'bb""b','cc"c'
'aa''a','b''
b''b',"c''cc"`

const csvWithComment = `aaa,bbb,ccc
; This is a comment.
"aaa",bbb,"ccc"
; This is another comment.
aaa,"b
bb","ccc"`

const csvWithCustomComment = `aaa,bbb,ccc
# This is a comment.
"aaa",bbb,"ccc"
# This is another comment.
aaa,"b
bb","ccc"`

const csvWithSpace = `aaa ,bbb, ccc
"aaa",  bbb  ,     "ccc"
 aaa, "b
  bb"  ,  "ccc"`

func TestScanner(t *testing.T) {
	var scanner = csv.NewScanner(strings.NewReader(csvStandard))
	_, rows, err := scanner.Scan()
	if err != nil {
		t.Error(err)
	}
	printRows(t, rows)
}

func TestScannerRFC4180(t *testing.T) {
	var scanner = csv.NewScanner(strings.NewReader(csvStandard), csv.RFC4180())
	_, rows, err := scanner.Scan()
	if err != nil {
		t.Error(err)
	}
	printRows(t, rows)
}

func TestScannerEscaped(t *testing.T) {
	var scanner = csv.NewScanner(strings.NewReader(csvEscaped))
	_, rows, err := scanner.Scan()
	if err != nil {
		t.Error(err)
	}
	printRows(t, rows)
}

func TestScannerWithHeader(t *testing.T) {
	var scanner = csv.NewScanner(strings.NewReader(csvWithHeader), csv.Header(true))
	header, rows, err := scanner.Scan()
	if err != nil {
		t.Error(err)
	}
	printHeader(t, header)
	printRows(t, rows)
}

func TestScannerWithCustomSeparator(t *testing.T) {
	var scanner = csv.NewScanner(strings.NewReader(csvWithCustomSeparator),
		csv.Header(true), csv.Separator('|'))
	header, rows, err := scanner.Scan()
	if err != nil {
		t.Error(err)
	}
	printHeader(t, header)
	printRows(t, rows)
}

func TestScannerWithSingleQuote(t *testing.T) {
	var scanner = csv.NewScanner(strings.NewReader(csvWithSingleQuote),
		csv.AllowSingleQuote(true))
	_, rows, err := scanner.Scan()
	if err != nil {
		t.Error(err)
	}
	printRows(t, rows)
}

func TestScannerWithComment(t *testing.T) {
	var scanner = csv.NewScanner(strings.NewReader(csvWithComment))
	_, rows, err := scanner.Scan()
	if err != nil {
		t.Error(err)
	}
	printRows(t, rows)
}

func TestScannerWithCustomComment(t *testing.T) {
	var scanner = csv.NewScanner(strings.NewReader(csvWithCustomComment), csv.Comment('#'))
	_, rows, err := scanner.Scan()
	if err != nil {
		t.Error(err)
	}
	printRows(t, rows)
}

func TestScannerWithSpace(t *testing.T) {
	var scanner = csv.NewScanner(strings.NewReader(csvWithSpace),
		csv.OmitLeadingSpace(true), csv.OmitTrailingSpace(true))
	_, rows, err := scanner.Scan()
	if err != nil {
		t.Error(err)
	}
	printRows(t, rows)
}

func printHeader(t *testing.T, header []string) {
	t.Logf("Header: [%s]\n", strings.Join(header, ", "))
}

func printRows(t *testing.T, rows [][]string) {
	for i, row := range rows {
		t.Logf("Row #%d: [%s]\n", i, strings.Join(row, ", "))
	}
}

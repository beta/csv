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
# This is a comment.
"aaa",bbb,"ccc"
# This is another comment.
aaa,"b
bb","ccc"`

const csvWithSpace = `aaa ,bbb, ccc
"aaa",  bbb  ,     "ccc"
 aaa, "b
  bb"  ,  "ccc"`

const csvWithEmptyField = `aaa,bbb,
"aaa",,"ccc"
,"b
  bb","ccc"`

const csvWithPrefixAndSuffix = `[Col A],[Col B],[Col C]
(aaa),(bbb),(ccc)
("aaa"),(bbb),("ccc")
(aaa),("b
bb"),("ccc")`

func TestScanner(t *testing.T) {
	_, rows, err := csv.Scan([]byte(csvStandard))
	if err != nil {
		t.Error(err)
	}
	printRows(t, rows)
}

func TestScannerRFC4180(t *testing.T) {
	_, rows, err := csv.Scan([]byte(csvStandard), csv.RFC4180())
	if err != nil {
		t.Error(err)
	}
	printRows(t, rows)
}

func TestScannerEscaped(t *testing.T) {
	_, rows, err := csv.Scan([]byte(csvEscaped))
	if err != nil {
		t.Error(err)
	}
	printRows(t, rows)
}

func TestScannerWithHeader(t *testing.T) {
	header, rows, err := csv.Scan([]byte(csvWithHeader), csv.Header(true))
	if err != nil {
		t.Error(err)
	}
	printHeader(t, header)
	printRows(t, rows)
}

func TestScannerWithCustomSeparator(t *testing.T) {
	header, rows, err := csv.Scan([]byte(csvWithCustomSeparator),
		csv.Header(true), csv.Separator('|'))
	if err != nil {
		t.Error(err)
	}
	printHeader(t, header)
	printRows(t, rows)
}

func TestScannerWithSingleQuote(t *testing.T) {
	_, rows, err := csv.Scan([]byte(csvWithSingleQuote),
		csv.AllowSingleQuote(true))
	if err != nil {
		t.Error(err)
	}
	printRows(t, rows)
}

func TestScannerWithComment(t *testing.T) {
	_, rows, err := csv.Scan([]byte(csvWithComment), csv.Comment('#'))
	if err != nil {
		t.Error(err)
	}
	printRows(t, rows)
}

func TestScannerWithSpace(t *testing.T) {
	_, rows, err := csv.Scan([]byte(csvWithSpace),
		csv.OmitLeadingSpace(true), csv.OmitTrailingSpace(true))
	if err != nil {
		t.Error(err)
	}
	printRows(t, rows)
}

func TestScannerWithEmptyField(t *testing.T) {
	_, rows, err := csv.Scan([]byte(csvWithEmptyField), csv.AllowEmptyField(true))
	if err != nil {
		t.Error(err)
	}
	printRows(t, rows)
}

func TestScannerWithPrefixAndSuffix(t *testing.T) {
	header, rows, err := csv.Scan([]byte(csvWithPrefixAndSuffix),
		csv.Header(true),
		csv.HeaderPrefix('['), csv.HeaderSuffix(']'),
		csv.FieldPrefix('('), csv.FieldSuffix(')'))
	if err != nil {
		t.Error(err)
	}
	printHeader(t, header)
	printRows(t, rows)
}

func printHeader(t *testing.T, header []string) {
	if len(header) > 0 {
		t.Logf("Header: [%s]\n", strings.Join(header, ", "))
	}
}

func printRows(t *testing.T, rows [][]string) {
	for i, row := range rows {
		t.Logf("Row #%d: [%s]\n", i, strings.Join(row, ", "))
	}
}

// Copyright (c) 2018 Beta Kuang
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package csv_test

import (
	"io"
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
	s, err := csv.NewScanner([]byte(csvStandard))
	if err != nil {
		t.Error(err)
		return
	}
	row, err := s.Scan()
	for err != io.EOF {
		if err != nil {
			t.Error(err)
			return
		}
		printRow(t, row)
		row, err = s.Scan()
	}
}

func TestScannerScanAll(t *testing.T) {
	s, err := csv.NewScanner([]byte(csvStandard))
	if err != nil {
		t.Error(err)
		return
	}
	rows, err := s.ScanAll()
	if err != nil {
		t.Error(err)
		return
	}
	printRows(t, rows)
}

func TestScannerRFC4180(t *testing.T) {
	s, err := csv.NewScanner([]byte(csvStandard), csv.RFC4180())
	if err != nil {
		t.Error(err)
		return
	}
	rows, err := s.ScanAll()
	if err != nil {
		t.Error(err)
		return
	}
	printRows(t, rows)
}

func TestScannerEscaped(t *testing.T) {
	s, err := csv.NewScanner([]byte(csvEscaped))
	if err != nil {
		t.Error(err)
		return
	}
	rows, err := s.ScanAll()
	if err != nil {
		t.Error(err)
		return
	}
	printRows(t, rows)
}

func TestScannerWithHeader(t *testing.T) {
	s, err := csv.NewScanner([]byte(csvWithHeader))
	if err != nil {
		t.Error(err)
		return
	}
	header, err := s.Scan()
	if err != nil {
		t.Error(err)
		return
	}
	printHeader(t, header)
	rows, err := s.ScanAll()
	if err != nil {
		t.Error(err)
		return
	}
	printRows(t, rows)
}

func TestScannerWithCustomSeparator(t *testing.T) {
	s, err := csv.NewScanner([]byte(csvWithCustomSeparator), csv.Separator('|'))
	if err != nil {
		t.Error(err)
		return
	}
	header, err := s.Scan()
	if err != nil {
		t.Error(err)
		return
	}
	printHeader(t, header)
	rows, err := s.ScanAll()
	if err != nil {
		t.Error(err)
		return
	}
	printRows(t, rows)
}

func TestScannerWithSingleQuote(t *testing.T) {
	s, err := csv.NewScanner([]byte(csvWithSingleQuote),
		csv.AllowSingleQuote(true))
	if err != nil {
		t.Error(err)
		return
	}
	rows, err := s.ScanAll()
	if err != nil {
		t.Error(err)
		return
	}
	printRows(t, rows)
}

func TestScannerWithComment(t *testing.T) {
	s, err := csv.NewScanner([]byte(csvWithComment), csv.Comment('#'))
	if err != nil {
		t.Error(err)
		return
	}
	rows, err := s.ScanAll()
	if err != nil {
		t.Error(err)
		return
	}
	printRows(t, rows)
}

func TestScannerWithSpace(t *testing.T) {
	s, err := csv.NewScanner([]byte(csvWithSpace),
		csv.OmitLeadingSpace(true), csv.OmitTrailingSpace(true))
	if err != nil {
		t.Error(err)
		return
	}
	rows, err := s.ScanAll()
	if err != nil {
		t.Error(err)
		return
	}
	printRows(t, rows)
}

func TestScannerWithEmptyField(t *testing.T) {
	s, err := csv.NewScanner([]byte(csvWithEmptyField), csv.AllowEmptyField(true))
	if err != nil {
		t.Error(err)
		return
	}
	rows, err := s.ScanAll()
	if err != nil {
		t.Error(err)
		return
	}
	printRows(t, rows)
}
func TestScannerWithPrefixAndSuffix(t *testing.T) {
	s, err := csv.NewScanner([]byte(csvWithPrefixAndSuffix))
	if err != nil {
		t.Error(err)
		return
	}

	s.Setting(csv.Prefix('['), csv.Suffix(']'))
	header, err := s.Scan()
	if err != nil {
		t.Error(err)
		return
	}
	printHeader(t, header)

	s.Setting(csv.Prefix('('), csv.Suffix(')'))
	rows, err := s.ScanAll()
	if err != nil {
		t.Error(err)
		return
	}
	printRows(t, rows)
}

func printHeader(t *testing.T, header []string) {
	if len(header) > 0 {
		t.Logf("Header: [%s]\n", strings.Join(header, ", "))
	}
}

func printRow(t *testing.T, row []string) {
	if len(row) > 0 {
		t.Logf("Row: [%s]\n", strings.Join(row, ", "))
	}
}

func printRows(t *testing.T, rows [][]string) {
	for i, row := range rows {
		t.Logf("Row #%d: [%s]\n", i, strings.Join(row, ", "))
	}
}

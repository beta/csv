// Copyright (c) 2018 Beta Kuang
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package csv_test

import (
	"testing"

	"github.com/beta/csv"
)

var records = [][]string{{"aaa", "bbb", "ccc"}, {"aaa", "b\nbb", "cc,c"}}

func TestGenerator(t *testing.T) {
	var g = csv.NewGenerator()
	var err error
	for _, record := range records {
		err = g.Write(record)
		if err != nil {
			t.Error(err)
			return
		}
	}
	data, err := g.Finish()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf(string(data))
}

func TestGeneratorWriteAll(t *testing.T) {
	var g = csv.NewGenerator()
	var err = g.WriteAll(records)
	if err != nil {
		t.Error(err)
		return
	}
	data, err := g.Finish()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf(string(data))
}

func TestGeneratorWithCustomSeparator(t *testing.T) {
	var g = csv.NewGenerator(csv.Separator('|'))
	var err = g.WriteAll(records)
	if err != nil {
		t.Error(err)
		return
	}
	data, err := g.Finish()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf(string(data))
}

func TestGeneratorWithPrefixSuffix(t *testing.T) {
	var g = csv.NewGenerator(csv.Prefix('('), csv.Suffix(')'))
	var err = g.WriteAll(records)
	if err != nil {
		t.Error(err)
		return
	}
	data, err := g.Finish()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf(string(data))
}

// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config_test

import (
	"testing"

	"github.com/aamcrae/config"
	"path/filepath"
	//    "reflect"
)

func TestString(t *testing.T) {
	s := `#
# Comment line 1

keyword=test
key2
key3=data1,data2,data3
`
	c, err := config.ParseString(s)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.GetSection(config.Global).GetEntries()) != 3 {
		t.Fatalf("Wrong number of config entries")
	}
	val := c.Get("keyword")
	if len(val) != 1 {
		t.Fatalf("wrong number of entries for keyword: %v", c)
	}
	if val[0].Lineno != 4 {
		t.Fatalf("wrong line number for 'keyword'")
	}
	if len(val[0].Tokens) != 1 {
		t.Fatalf("Wrong number of tokens for 'keyword'")
	}
	if val[0].Tokens[0] != "test" {
		t.Fatalf("Wrong token value for 'keyword'")
	}
	val = c.Get("key2")
	if len(val) != 1 {
		t.Fatalf("wrong number of entries for 'key2'")
	}
	if val[0].Lineno != 5 {
		t.Fatalf("wrong line number for 'key2'")
	}
	if len(val[0].Tokens) != 0 {
		t.Fatalf("Wrong number of tokens for 'key2'")
	}
	val = c.Get("key3")
	if len(val) != 1 {
		t.Fatalf("wrong number of entries for 'key3'")
	}
	if val[0].Lineno != 6 {
		t.Fatalf("wrong line number for 'key3'")
	}
	if len(val[0].Tokens) != 3 {
		t.Fatalf("Wrong number of tokens for 'key3'")
	}
}

func TestFile(t *testing.T) {
	c, err := config.ParseFile(filepath.Join("testdata", "f1"))
	if err != nil {
		t.Fatalf("File read for f1 failed: %v", err)
	}
	l := len(c.GetSection(config.Global).GetEntries())
	if l != 3 {
		t.Fatalf("TestFile: wrong number of entries: %d", l)
	}
	l = len(c.GetSection("section1").GetEntries())
	if l != 2 {
		t.Fatalf("TestFile: wrong number of entries for 'section1': %d", l)
	}
	val := c.Get("key1")
	if len(val) != 1 {
		t.Fatalf("wrong number of entries for 'key1'")
	}
	exp := [...]string{"1", "2", "3", "4"}
	if len(val[0].Tokens) != len(exp) {
		t.Fatalf("TestFile: Wrong number of tokens for 'key1'")
	}
	for i, v := range exp {
		if val[0].Tokens[i] != v {
			t.Fatalf("TestFile: Wrong token %d for 'key1', expected %s, got %s", i, v, val[0].Tokens[i])
		}
	}
	s1 := c.GetSection("section1")
	val = s1.Get("key1")
	if len(val) != 1 {
		t.Fatalf("wrong number of entries for 'section1' 'key1'")
	}
	exp1 := [...]string{"x", "y", "z"}
	if len(val[0].Tokens) != len(exp1) {
		t.Fatalf("TestFile: Wrong number of tokens for 'section1' 'key1'")
	}
	for i, v := range exp1 {
		if val[0].Tokens[i] != v {
			t.Fatalf("TestFile: Wrong token %d for 'section1' 'key1', expected %s, got %s", i, v, val[0].Tokens[i])
		}
	}
}

func TestMultiFile(t *testing.T) {
	p1 := filepath.Join("testdata", "f1")
	p2 := filepath.Join("testdata", "f2")
	c, err := config.ParseFiles(true, []string{p1, p2})
	if err != nil {
		t.Fatalf("File read for f1/f2 failed: %v", err)
	}
	l := len(c.GetSection(config.Global).GetEntries())
	if l != 6 {
		t.Fatalf("TestFiles: wrong number of entries: %d", l)
	}
	val := c.Get("key3")
	if len(val) != 2 {
		t.Fatalf("wrong number of entries for 'key1'")
	}
	if len(val[0].Tokens) != 1 || len(val[1].Tokens) != 1 {
		t.Fatalf("TestFiles: Wrong number of tokens for 'key3'")
	}
	if val[0].Tokens[0] != "abc" || val[1].Tokens[0] != "xyz" {
		t.Fatalf("TestFiles: Incorrect values for 'key3'")
	}
}

func TestApi(t *testing.T) {
	c, err := config.ParseFile(filepath.Join("testdata", "f1"))
	if err != nil {
		t.Fatalf("TestApi: File read for f1 failed: %v", err)
	}
	v := c.Get("key1")
	if len(v) != 1 {
		t.Fatalf("wrong number of entries for 'key1'")
	}
	if len(v[0].Tokens) != 4 {
		t.Fatalf("TestApi: Get returns wrong number of tokens: %v", v[0].Tokens)
	}
	exp := []string{"1", "2", "3", "4"}
	for i, tok := range exp {
		if v[0].Tokens[i] != tok {
			t.Fatalf("TestApi: Get returns wrong tokens : %v", v[0].Tokens)
		}
	}
	strs := c.Missing([]string{"key2", "key1", "key5"})
	if len(strs) != 1 {
		t.Fatalf("TestApi: Missing returns wrong number: %v", strs)
	}
	if strs[0] != "key5" {
		t.Fatalf("TestApi: Missing returns wrong string: %v", strs)
	}
}

func TestMerge(t *testing.T) {
	c, err := config.ParseFile(filepath.Join("testdata", "f1"))
	if err != nil {
		t.Fatalf("TestMerge: File read for f1 failed: %v", err)
	}
	c1, err := config.ParseFile(filepath.Join("testdata", "f2"))
	if err != nil {
		t.Fatalf("TestMerge: File read for f2 failed: %v", err)
	}
	p1 := filepath.Join("testdata", "f1")
	p2 := filepath.Join("testdata", "f2")
	comp, err := config.ParseFiles(true, []string{p1, p2})
	if err != nil {
		t.Fatalf("TestMerge: File read for f1/f2 failed: %v", err)
	}
	c.Merge(c1)
	// A bit tricky to compare, since the values are pointers.
	l := len(c.GetSection(config.Global).GetEntries())
	if l != len(comp.GetSection(config.Global).GetEntries()) {
		t.Fatalf("TestMerge: lengths are different: %v %v", c1, comp)
	}
}

func TestSection(t *testing.T) {
	c, err := config.ParseFile(filepath.Join("testdata", "f3"))
	if err != nil {
		t.Fatalf("TestMerge: File read for f3 failed: %v", err)
	}
	// There should be 2 sections named 'section'
	s := c.GetSections("section")
	if len(s) != 2 {
		t.Fatalf("section 'section' in 'f3': exp 2, got %d", len(s))
	}
	v := s[1].Get("key2")
	if len(v) != 1 {
		t.Fatalf("missing 'key1' in f3 section #2")
	}
}

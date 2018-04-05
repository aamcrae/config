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
    if len(c.Entries) != 3 {
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
    if len(c.Entries) != 3 {
        t.Fatalf("TestFile: wrong number of entries: %d", len(c.Entries))
    }
    val := c.Get("key1")
    if len(val) != 1 {
        t.Fatalf("wrong number of entries for 'key1'")
    }
    if len(val[0].Tokens) != 4 {
        t.Fatalf("TestFile: Wrong number of tokens for 'key1'")
    }
}

func TestMultiFile(t *testing.T) {
    p1 := filepath.Join("testdata", "f1")
    p2 := filepath.Join("testdata", "f2")
    c, err := config.ParseFiles(true, []string{p1, p2})
    if err != nil {
        t.Fatalf("File read for f1/f2 failed: %v", err)
    }
    if len(c.Entries) != 6 {
        t.Fatalf("TestFiles: wrong number of entries: %d", len(c.Entries))
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
    if len(c.Entries) != len(comp.Entries) {
        t.Fatalf("TestMerge: lengths are different: %v %v", c1, comp)
    }
}

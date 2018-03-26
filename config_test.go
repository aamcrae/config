package config_test

import (
    "testing"

    "github.com/aamcrae/config"
    "path/filepath"
    "reflect"
)

func TestString(t *testing.T) {
    s := `#
# Comment line 1

keyword = test
key2
key3=data1,data2,data3
`
    c, err := config.ParseString(s)
    if err != nil {
        t.Fatal(err)
    }
    if len(c) != 3 {
        t.Fatalf("Wrong number of config entries")
    }
    val, ok := c["keyword"]
    if !ok {
        t.Fatalf("keyword 'keyword' has not been found")
    }
    if val.Lineno != 4 {
        t.Fatalf("wrong line number for 'keyword'")
    }
    if len(val.Tokens) != 1 {
        t.Fatalf("Wrong number of tokens for 'keyword'")
    }
    if val.Tokens[0] != "test" {
        t.Fatalf("Wrong token value for 'keyword'")
    }
    val, ok = c["key2"]
    if !ok {
        t.Fatalf("keyword 'key2' has not been found")
    }
    if val.Lineno != 5 {
        t.Fatalf("wrong line number for 'key2'")
    }
    if len(val.Tokens) != 0 {
        t.Fatalf("Wrong number of tokens for 'key2'")
    }
    val, ok = c["key3"]
    if !ok {
        t.Fatalf("keyword 'key3' has not been found")
    }
    if val.Lineno != 6 {
        t.Fatalf("wrong line number for 'key3'")
    }
    if len(val.Tokens) != 3 {
        t.Fatalf("Wrong number of tokens for 'key3'")
    }
}

func TestFile(t *testing.T) {
    c, err := config.ParseFile(filepath.Join("testdata", "f1"))
    if err != nil {
        t.Fatalf("File read for f1 failed: %v", err)
    }
    if len(c) != 3 {
        t.Fatalf("TestFile: wrong number of entries: %d", len(c))
    }
    val, ok := c["key1"]
    if !ok {
        t.Fatalf("TestFile: keyword 'key1' has not been found")
    }
    if len(val.Tokens) != 4 {
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
    if len(c) != 5 {
        t.Fatalf("TestFiles: wrong number of entries: %d", len(c))
    }
    val, ok := c["key3"]
    if !ok {
        t.Fatalf("TestFiles: keyword 'key3' has not been found")
    }
    if len(val.Tokens) != 1 {
        t.Fatalf("TestFiles: Wrong number of tokens for 'key3'")
    }
    if val.Tokens[0] != "xyz" {
        t.Fatalf("TestFiles: Incorrect value for 'key3'")
    }
}

func TestApi(t *testing.T) {
    c, err := config.ParseFile(filepath.Join("testdata", "f1"))
    if err != nil {
        t.Fatalf("TestApi: File read for f1 failed: %v", err)
    }
    v, ok := c.Get("key1")
    if !ok {
        t.Fatalf("TestApi: Get failed on key 'key1'")
    }
    if len(v.Tokens) != 4 {
        t.Fatalf("TestApi: Get returns wrong number of tokens: %v", v.Tokens)
    }
    exp := []string{"1", "2", "3", "4"}
    for i, tok := range exp {
        if v.Tokens[i] != tok {
            t.Fatalf("TestApi: Get returns wrong tokens : %v", v.Tokens)
        }
    }
    vals := c.GetN([]string{"key1", "key2"})
    if len(vals) != 2 {
        t.Fatalf("TestApi: GetN returns wrong number, vals: %v", vals)
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
    if len(c) != len(comp) {
        t.Fatalf("TestMerge: lengths are different: %v %v", c1, comp)
    }
    for k, v := range c {
        if vc, ok := comp[k]; !ok {
            t.Fatalf("TestMerge: %v missing", k)
        } else if !reflect.DeepEqual(*v, *vc) {
            t.Fatalf("TestMerge: %v values different: %v %v", k, *v, *vc)
        }
    }
}

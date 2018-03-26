package config_test

import (
    "testing"

    "path/filepath"
    "github.com/aamcrae/config"
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
    path := filepath.Join("testdata", "f1")
    c, err := config.ParseFile(path)
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

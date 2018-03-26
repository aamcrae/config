package config

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

// Value, one for each keyword.
type Value struct {
    Filename string
    Lineno int
    Line string
    Tokens []string
}

type Config map[string]*Value

func (c *Config) Merge(c1 Config) {
    for k, v := range c1 {
        (*c)[k] = v
    }
}

func (c *Config) Get(k string) ([]string, bool) {
    if v, ok := (*c)[k]; ok {
        return v.Tokens, true
    }
    return nil, false
}

func (c *Config) ParseFile(file string) error {
    return parseOneFile(file, *c)
}

func (c *Config) GetN(strs []string) []*Value {
    var values []*Value
    for _, s := range strs {
        if v, ok := (*c)[s]; ok {
            values = append(values, v)
        }
    }
    return values
}

func ParseFiles(optional bool, files []string) (Config, error) {
    config := make(Config)
    for _, f := range files {
        if err := parseOneFile(f, config); err != nil {
            if !optional {
                return config, fmt.Errorf("%s: %v", f, err)
            }
        }
    }
    return config, nil
}

func ParseFile(file string) (Config, error) {
    config := make(Config)
    return config, parseOneFile(file, config)
}

func parseOneFile(file string, config Config) error {
    f, err := os.Open(file)
    if err != nil {
        return err
    }
    defer f.Close()
    return parse(file, bufio.NewReader(f), config)
}

func ParseString(s string) (Config, error) {
    config := make(Config)
    return config, parse("internal", bufio.NewReader(strings.NewReader(s)), config)
}

// parse reads the input and places key/value pairs in Config.
// Comments are marked as '#' at the start of the line.
// The lines are expected in the format:
//   keyword[ [ = ] tokens]
// Tokens are delimited by space, comma, tabs or '='
// Duplicate keywords are silently overwritten.
func parse(source string, r *bufio.Reader, config Config) error {
    lineno := 0
    scanner := bufio.NewScanner(r)
    for scanner.Scan() {
        lineno++
        l := strings.TrimSpace(scanner.Text())
        if len(l) == 0 || l[0:1] == "#" {
            continue
        }
        tok := strings.FieldsFunc(l, delimiters)
        if len(tok) == 0 {
            continue
        }
        v := new(Value)
        v.Filename = source
        v.Lineno = lineno
        v.Line = l
        v.Tokens = tok[1:]
        config[tok[0]] = v
    }
    if scanner.Err() != nil {
        return fmt.Errorf("%s: line %d: %v", source, lineno, scanner.Err())
    }
    return nil
}

func delimiters(r rune) bool {
    return r == ' ' || r == '=' || r == '\t' || r == ','
}

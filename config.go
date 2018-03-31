package config

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

// Value, one for each keyword.
type Value struct {
    Keyword string
    Filename string
    Lineno int
    Line string
    Tokens []string
    index int
}

type Config struct {
    m map[string]*Value
    Values []*Value
}

func (c *Config) Merge(c1 *Config) {
    for _, v := range c1.Values {
        c.addValue(v)
    }
}

func (c *Config) GetTokens(k string) ([]string, bool) {
    if v, ok := c.m[k]; ok {
        return v.Tokens, true
    }
    return nil, false
}

func (c *Config) Get(k string) (*Value, bool) {
    if v, ok := c.m[k]; ok {
        return v, true
    }
    return nil, false
}

func (c *Config) ParseFile(file string) error {
    return c.parseOneFile(file)
}

func (c *Config) GetN(strs []string) []*Value {
    var values []*Value
    for _, s := range strs {
        if v, ok := c.m[s]; ok {
            values = append(values, v)
        }
    }
    return values
}

// Return strings not present in config.
func (c *Config) Missing(strs []string) []string {
    var missing []string
    for _, s := range strs {
        if _, ok := c.m[s]; !ok {
            missing = append(missing, s)
        }
    }
    return missing
}

func ParseFiles(optional bool, files []string) (*Config, error) {
    config := &Config{map[string]*Value{}, []*Value{}}
    for _, f := range files {
        if err := config.parseOneFile(f); err != nil {
            if !optional {
                return config, fmt.Errorf("%s: %v", f, err)
            }
        }
    }
    return config, nil
}

func ParseFile(file string) (*Config, error) {
    config := &Config{map[string]*Value{}, []*Value{}}
    return config, config.parseOneFile(file)
}

func (config *Config) parseOneFile(file string) error {
    f, err := os.Open(file)
    if err != nil {
        return err
    }
    defer f.Close()
    return config.parse(file, bufio.NewReader(f))
}

func ParseString(s string) (*Config, error) {
    config := &Config{map[string]*Value{}, []*Value{}}
    return config, config.parse("internal", bufio.NewReader(strings.NewReader(s)))
}

// parse reads the input and places key/value pairs in Config.
// Comments are marked as '#' at the start of the line.
// The lines are expected in the format:
//   keyword[ [ = ] tokens]
// Tokens are delimited by space, comma, tabs or '='
// Duplicate keywords are silently overwritten.
func (config *Config) parse(source string, r *bufio.Reader) error {
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
        config.addValue(&Value{tok[0], source, lineno, l, tok[1:], 0})
    }
    if scanner.Err() != nil {
        return fmt.Errorf("%s: line %d: %v", source, lineno, scanner.Err())
    }
    return nil
}

func (config *Config) addValue(v *Value) {
    entry, ok := config.m[v.Keyword]
    if ok {
        // Entry already exists, overwrite the existing one.
        v.index = entry.index
        config.Values[v.index] = v
    } else {
        // New entry, add to the end of the list.
        v.index = len(config.Values)
        config.Values = append(config.Values, v)
    }
    config.m[v.Keyword] = v
}

func delimiters(r rune) bool {
    return r == ' ' || r == '=' || r == '\t' || r == ','
}

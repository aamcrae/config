package config

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

// Entry, one for each keyword.
type Entry struct {
    Keyword string
    Filename string
    Lineno int
    Line string
    Tokens []string
}

type Config struct {
    m map[string][]*Entry
    Entries []*Entry
}

var delimiters map[rune]struct{} = map[rune]struct{}{ '=':struct{}{}, ',':struct{}{} }

func SetDelimiters(d string) {
    delimiters = map[rune]struct{}{}
    for _, r := range d {
        delimiters[r] = struct{}{}
    }
}

func (c *Config) Merge(c1 *Config) {
    for _, v := range c1.Entries {
        c.addEntry(v)
    }
}

func (c *Config) Get(k string) []*Entry {
    if v, ok := c.m[k]; ok {
        return v
    }
    return []*Entry{}
}

func (c *Config) GetArg(k string) (string, error) {
    if v, ok := c.m[k]; ok {
        if len(v) != 1 {
            return "", fmt.Errorf("Illegal config for '%s'", k)
        }
        if len(v[0].Tokens) != 1 {
            return "", fmt.Errorf("Illegal arguments for '%s'", k)
        }
        return v[0].Tokens[0], nil
    }
    return "", fmt.Errorf("Missing keyword: %s", k)
}

func (c *Config) ParseFile(file string) error {
    return c.parseOneFile(file)
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
    config := &Config{map[string][]*Entry{}, []*Entry{}}
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
    config := &Config{map[string][]*Entry{}, []*Entry{}}
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
    config := &Config{map[string][]*Entry{}, []*Entry{}}
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
        tok := strings.FieldsFunc(l, checkDelimiter)
        if len(tok) == 0 {
            continue
        }
        config.addEntry(&Entry{tok[0], source, lineno, l, tok[1:]})
    }
    if scanner.Err() != nil {
        return fmt.Errorf("%s: line %d: %v", source, lineno, scanner.Err())
    }
    return nil
}

func (config *Config) addEntry(v *Entry) {
    config.Entries = append(config.Entries, v)
    entry, ok := config.m[v.Keyword]
    if !ok {
        entry = []*Entry{}
    }
    config.m[v.Keyword] = append(entry, v)
}

func checkDelimiter(r rune) bool {
    _, found := delimiters[r]
    return found
}

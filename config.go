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

// The entries associated with a single section.
// The default section is named 'global'.
type Section struct {
    m map[string][]*Entry
    entries []*Entry
}

type Config struct {
	sections map[string]*Section
}

const Global = "global"

var delimiters map[rune]struct{} = map[rune]struct{}{ '=':struct{}{}, ',':struct{}{} }

func SetDelimiters(d string) {
    delimiters = map[rune]struct{}{}
    for _, r := range d {
        delimiters[r] = struct{}{}
    }
}

func (c *Config) Merge(c1 *Config) {
    for s, v := range c1.sections {
		sect := c.getSection(s, true)
		for _, e := range v.entries {
			sect.addEntry(e)
		}
    }
}

func (c *Config) Has(k string) bool {
	return c.GetSection(Global).Has(k)
}

func (c *Config) Get(k string) []*Entry {
	return c.GetSection(Global).Get(k)
}

func (c *Config) GetArg(k string) (string, error) {
	return c.GetSection(Global).GetArg(k)
}

func (config *Config) GetSection(name string) *Section {
	return config.getSection(name, false)
}

func (config *Config) getSection(name string, create bool) *Section {
	if name == "" {
		name = Global
	}
	s, ok := config.sections[name]
	if !ok {
		if !create {
			return nil
		}
		s = &Section{map[string][]*Entry{}, []*Entry{}}
		config.sections[name] = s
	}
	return s
}


func (c *Config) ParseFile(file string) error {
    return c.parseOneFile(file)
}

// Return strings not present in config.
func (c *Config) Missing(strs []string) []string {
    var missing []string
    for _, s := range strs {
        if _, ok := c.sections[Global].m[s]; !ok {
            missing = append(missing, s)
        }
    }
    return missing
}

func (s *Section) Has(k string) bool {
    _, ok := s.m[k]
	return ok
}

func (s *Section) Get(k string) []*Entry {
    if v, ok := s.m[k]; ok {
        return v
    }
    return []*Entry{}
}

func (s *Section) GetArg(k string) (string, error) {
    if v, ok := s.m[k]; ok {
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

func (s *Section) GetEntries() []*Entry {
	return s.entries
}

func ParseFiles(optional bool, files []string) (*Config, error) {
	config := newConfig()
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
	config := newConfig()
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
	config := newConfig()
    return config, config.parse("internal", bufio.NewReader(strings.NewReader(s)))
}

// newConfig creates a new Config structure.
func newConfig() *Config {
	c := new(Config)
	c.sections = make(map[string]*Section)
	c.getSection(Global, true)
	return c
}

// parse reads the input and places key/value pairs in Config.
// Comments are marked as '#' at the start of the line.
// Separate sections are marked as:
//  [section-name]
// The lines are expected in the format:
//   keyword[ [ = ] tokens]
// Tokens are delimited by space, comma, tabs or '='
// Duplicate keywords are silently overwritten.
func (config *Config) parse(source string, r *bufio.Reader) error {
	sect := config.sections[Global]
    lineno := 0
    scanner := bufio.NewScanner(r)
    for scanner.Scan() {
        lineno++
        l := strings.TrimSpace(scanner.Text())
        if len(l) == 0 || l[0:1] == "#" {
            continue
        }
		ln := len(l)
		// Check for new section.
		if ln > 2 && l[0] == '[' && l[ln-1] == ']' {
			sect = config.getSection(l[1:ln-1], true)
			continue
		}
        tok := strings.FieldsFunc(l, checkDelimiter)
        if len(tok) == 0 {
            continue
        }
        sect.addEntry(&Entry{tok[0], source, lineno, l, tok[1:]})
    }
    if scanner.Err() != nil {
        return fmt.Errorf("%s: line %d: %v", source, lineno, scanner.Err())
    }
    return nil
}

func (sect *Section) addEntry(v *Entry) {
    sect.entries = append(sect.entries, v)
    entry, ok := sect.m[v.Keyword]
    if !ok {
        entry = []*Entry{}
    }
    sect.m[v.Keyword] = append(entry, v)
}

func checkDelimiter(r rune) bool {
    _, found := delimiters[r]
    return found
}

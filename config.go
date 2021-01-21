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

// Interface for accessing configuration tokens.
type Conf interface {
	Has(string) bool
	Get(string) []*Entry
	GetArg(string) (string, error)
}

// Entry, one for each keyword.
type Entry struct {
	Keyword  string
	Filename string
	Lineno   int
	Line     string
	Args	 string		// String following keyword
	Tokens   []string
}

// The entries associated with a single section.
// The default section is named 'global'.
type Section struct {
	Name string
	m       map[string][]*Entry
	entries []*Entry
}

type Config struct {
	sections map[string][]*Section
}

const Global = "global"

var delimiters map[rune]struct{} = map[rune]struct{}{'=': struct{}{}, ',': struct{}{}}

func SetDelimiters(d string) {
	delimiters = map[rune]struct{}{}
	for _, r := range d {
		delimiters[r] = struct{}{}
	}
}

func (c *Config) Merge(c1 *Config) {
	for s, v := range c1.sections {
		sect := c.addSection(s)
		for _, sl := range v {
			for _, e := range sl.entries {
				sect.addEntry(e)
			}
		}
	}
}

func (c *Config) Has(k string) bool {
	return c.GetSection(Global).Has(k)
}

func (c *Config) Parse(k string, f string, a ...interface{}) (int, error) {
	e := c.GetSection(Global).Get(k)
	if len(e) != 1 {
		return 0, fmt.Errorf("%s: invalid keyword(s)", k)
	}
	return fmt.Sscanf(e[0].Args, f, a...)
}

func (c *Config) Get(k string) []*Entry {
	return c.GetSection(Global).Get(k)
}

func (c *Config) GetArg(k string) (string, error) {
	return c.GetSection(Global).GetArg(k)
}

func (config *Config) GetSections(name string) []*Section {
	return config.getSections(name)
}

// Get the first section.
func (config *Config) GetSection(name string) *Section {
	s := config.getSections(name)
	if len(s) != 0 {
		return s[0]
	}
	return nil
}

// getSections retrieves the named slice of sections.
func (config *Config) getSections(name string) []*Section {
	if name == "" {
		name = Global
	}
	return config.sections[name]
}

// addSection adds the named section. Global is treated specially, where
// only a single Section is kept.
func (config *Config) addSection(name string) *Section {
	if name == "" {
		name = Global
	}
	s := config.sections[name]
	if name == Global && len(s) == 1 {
		return s[0]
	}
	sect := &Section{name, map[string][]*Entry{}, []*Entry{}}
	config.sections[name] = append(s, sect)
	return sect
}

func (c *Config) ParseFile(file string) error {
	return c.parseOneFile(file)
}

// Return strings not present in config.
func (c *Config) Missing(strs []string) []string {
	var missing []string
	for _, s := range strs {
		if _, ok := c.sections[Global][0].m[s]; !ok {
			missing = append(missing, s)
		}
	}
	return missing
}

func (s *Section) Has(k string) bool {
	_, ok := s.m[k]
	return ok
}

func (s *Section) Parse(k string, f string, a ...interface{}) (int, error) {
	if v, ok := s.m[k]; ok {
		if len(v) != 1 {
			return 0, fmt.Errorf("%s: invalid keyword(s)", k)
		}
		return fmt.Sscanf(v[0].Args, f, a...)
	}
	return 0, fmt.Errorf("%s: keyword not found", k)
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
	c.sections = make(map[string][]*Section)
	c.addSection(Global)
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
	sect := config.GetSection(Global)
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
			sect = config.addSection(l[1:ln-1])
			continue
		}
		tok := strings.FieldsFunc(l, checkDelimiter)
		if len(tok) == 0 {
			continue
		}
		var sb strings.Builder
		for i, s := range tok[1:] {
			if i != 0 {
				sb.WriteRune(',')
			}
			sb.WriteString(s)
		}
		sect.addEntry(&Entry{tok[0], source, lineno, l, sb.String(), tok[1:]})
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

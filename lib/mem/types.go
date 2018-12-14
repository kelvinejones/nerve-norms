package mem

import (
	"regexp"
	"strings"
)

type XY struct {
	X float64
	Y float64
}

type XYZ struct {
	X float64
	Y float64
	Z float64
}

type LineParser interface {
	LinePrefix() string

	// ParseRegex provides a regexp used to split a line.
	ParseRegex() *regexp.Regexp

	// ParseLine parses a line that was split by the Regexp.
	// `err` might be non-nil even if `keepParsing` is true; it's not a terminating error.
	ParseLine([]string) error
}

type section interface {
	Header() []string
	LineParser
	Parse(reader *Reader) error
}

// sectionHeaderMatches returns true if one of the section's headers matches this string.
func sectionHeaderMatches(sec section, str string) bool {
	for _, hd := range sec.Header() {
		if strings.Contains(str, hd) {
			return true
		}
	}
	return false
}

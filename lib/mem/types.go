package mem

import (
	"regexp"
)

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

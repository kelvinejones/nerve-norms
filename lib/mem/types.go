package mem

import "regexp"

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
	ParseRegex() *regexp.Regexp
	ParseLine([]string) error
}

type section interface {
	Header() string
	LineParser
	Parse(reader *Reader) error
}

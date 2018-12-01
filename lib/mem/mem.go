package mem

import "time"

type Sex int

const (
	UnknownSex Sex = iota
	FemaleSex
	MaleSex
)

type Header struct {
	File      string
	Name      string
	Protocol  string
	Date      time.Time
	StartTime time.Time // TODO get rid of this field; merge into Date
	Age       int
	Sex
	Temperature   float32
	SRSites       string
	NormalControl bool
	Operator      string
	Comment       string
}

type StimResponse struct {
	MaxCmap float32
}

type Mem struct {
	Header
	StimResponse
}

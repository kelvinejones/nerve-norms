package mem

import "time"

type Sex int

const (
	MaleSex Sex = iota
	FemaleSex
	OtherSex
)

type MemHeader struct {
	File     string
	Name     string
	Protocol string
	Date     time.Time
	Age      int
	Sex
	Temperature   float32
	SRSites       string
	NormalControl bool
	Operator      string
	Comment       string
}

type Mem struct {
	MemHeader
}

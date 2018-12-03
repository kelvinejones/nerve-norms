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
	Temperature   float64
	SRSites       string
	NormalControl bool
	Operator      string
	Comment       string
}

type StimResponse struct {
	MaxCmap float64
	Values  []XY
}

type ChargeDuration struct {
	Values []XYZ
}

type ThresholdElectrotonusGroup struct {
	Sets []ThresholdElectrotonusSet
}

type ThresholdElectrotonusSet struct {
	Values []XYZ
}

type Mem struct {
	Header
	StimResponse
	ChargeDuration
	ThresholdElectrotonusGroup
}

type XY struct {
	X float64
	Y float64
}

type XYZ struct {
	X float64
	Y float64
	Z float64
}

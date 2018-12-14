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

type MaxCmap struct {
	Val   float64
	Time  float64
	Units byte
}

type MaxCmaps []MaxCmap

type StimResponse struct {
	MaxCmaps
	Values    []XY
	ValueType string
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

type RecoveryCycle struct {
	Values []XY
}

type ThresholdIV struct {
	Values []XY
}

type ExcitabilityVariables struct {
	Values          map[string]float64
	Program         string
	ThresholdMethod int
	SRMethod        int
}

type Mem struct {
	Header
	StimResponse
	ChargeDuration
	ThresholdElectrotonusGroup
	RecoveryCycle
	ThresholdIV
	ExcitabilityVariables
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

type section interface {
	Header() string
	Parser
}

func (section StimResponse) Header() string {
	return "STIMULUS-RESPONSE DATA"
}

func (section ChargeDuration) Header() string {
	return "CHARGE DURATION DATA"
}

func (section ThresholdElectrotonusGroup) Header() string {
	return "THRESHOLD ELECTROTONUS DATA"
}

func (section RecoveryCycle) Header() string {
	return "RECOVERY CYCLE DATA"
}

func (section ThresholdIV) Header() string {
	return "THRESHOLD I/V DATA"
}

func (section ExcitabilityVariables) Header() string {
	return "DERIVED EXCITABILITY VARIABLES"
}

type ExtraVariables struct {
	*ExcitabilityVariables
}

func (section ExtraVariables) Header() string {
	return "EXTRA VARIABLES"
}

// This section has not been implemented, so skip it
type StrengthDuration struct{}

func (section StrengthDuration) Header() string {
	return "STRENGTH-DURATION DATA"
}

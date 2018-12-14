package mem

import (
	"fmt"
	"time"
)

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

func (header Header) String() string {
	return "Header{File{\"" + header.File + "\"}, Name{\"" + header.Name + "\"} }"
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

func (sr StimResponse) String() string {
	return fmt.Sprintf("StimResponse{%d MaxCmaps, %d values}", len(sr.MaxCmaps), len(sr.Values))
}

type ChargeDuration struct {
	Values []XYZ
}

func (cd ChargeDuration) String() string {
	return fmt.Sprintf("ChargeDuration{%d values}", len(cd.Values))
}

type ThresholdElectrotonusGroup struct {
	Sets []ThresholdElectrotonusSet
}

func (teg ThresholdElectrotonusGroup) String() string {
	str := "ThresholdElectrotonusGroup{"
	for _, tes := range teg.Sets {
		str += tes.String() + ","
	}
	str += "}"
	return str
}

type ThresholdElectrotonusSet struct {
	Values []XYZ
}

func (tes ThresholdElectrotonusSet) String() string {
	return fmt.Sprintf("ThresholdElectrotonusSet{%d values}", len(tes.Values))
}

type RecoveryCycle struct {
	Values []XY
}

func (rc RecoveryCycle) String() string {
	return fmt.Sprintf("RecoveryCycle{%d values}", len(rc.Values))
}

type ThresholdIV struct {
	Values []XY
}

func (tiv ThresholdIV) String() string {
	return fmt.Sprintf("ThresholdIV{%d values}", len(tiv.Values))
}

type ExcitabilityVariables struct {
	Values          map[string]float64
	Program         string
	ThresholdMethod int
	SRMethod        int
}

func (ev ExcitabilityVariables) String() string {
	return fmt.Sprintf("ExcitabilityVariables{%d values}", len(ev.Values))
}

type Mem struct {
	Header
	StimResponse
	ChargeDuration
	ThresholdElectrotonusGroup
	RecoveryCycle
	ThresholdIV
	ExcitabilityVariables
	StrengthDuration
}

func (mem Mem) String() string {
	str := "Mem{\n"
	str += "\t" + mem.Header.String() + ",\n"
	str += "\t" + mem.StimResponse.String() + ",\n"
	str += "\t" + mem.ChargeDuration.String() + ",\n"
	str += "\t" + mem.ThresholdElectrotonusGroup.String() + ",\n"
	str += "\t" + mem.RecoveryCycle.String() + ",\n"
	str += "\t" + mem.ThresholdIV.String() + ",\n"
	str += "\t" + mem.ExcitabilityVariables.String() + ",\n"
	str += "\t" + mem.StrengthDuration.String() + ",\n"
	str += "}"
	return str
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

func (sd StrengthDuration) String() string {
	return fmt.Sprintf("StrengthDuration{Import not implemented}")
}

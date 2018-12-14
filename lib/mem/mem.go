package mem

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

package mem

type ThresholdIV struct{ LabTab }

var IVCurrent = Column([]float64{50, 40, 30, 20, 10, 0, -10, -20, -30, -40, -50, -60, -70, -80, -90, -100})

func IVLabelledTable(mem *Mem) LabelledTable {
	return &mem.Sections["IV"].(*ThresholdIV).LabTab
}

func newIV() *ThresholdIV {
	return &ThresholdIV{LabTab{
		section:    "THRESHOLD I/V",
		altSection: "THESHOLD I/V",
		xname:      "Current (%)",
		yname:      "Threshold redn. (%)",
		altYname:   "Threshold change (%)",
		xcol:       IVCurrent,
		precision:  0.01,
	}}
}

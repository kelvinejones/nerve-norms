package mem

type ChargeDuration struct{ LabTab }

var CDDuration = Column([]float64{0.2, 0.4, 0.6, 0.8, 1})

func CDLabelledTable(mem *Mem) LabelledTable {
	return &mem.Sections["CD"].(*ChargeDuration).LabTab
}

func newCD() *ChargeDuration {
	return &ChargeDuration{LabTab{
		section:  "CHARGE DURATION",
		xname:    "Duration (ms)",
		altXname: "Current (%)",
		yname:    "Threshold charge (mA.mS)",
		altYname: "Threshold change (%)",
		xcol:     CDDuration,
		altImportFunc: func(cd *LabTab) {
			for i := range cd.ycol {
				cd.ycol[i] *= cd.xcol[i]
			}
		},
		precision: 0.0000001,
	}}
}

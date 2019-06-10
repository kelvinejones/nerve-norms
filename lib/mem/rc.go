package mem

type RecoveryCycle struct{ LabTab }

var RCInterval = Column([]float64{2, 2.5, 3.2, 4, 5, 6.3, 7.9, 10, 13, 18, 24, 32, 42, 56, 75, 100, 140, 200})

func newRC() *RecoveryCycle {
	rc := &RecoveryCycle{LabTab{
		section:   "RECOVERY CYCLE",
		xname:     "Interval (ms)",
		yname:     "Threshold change (%)",
		xcol:      RCInterval,
		precision: 0.000001,
		logScale:  true,
	}}
	rc.LabTab.postImputeAction = func() {
		maxRC0 := 200.0
		maxRC1 := 100.0
		if rc.ycol[0] < maxRC0 && rc.ycol[1] < maxRC1 {
			return
		}
		if rc.wasimp == nil {
			rc.wasimp = make(Column, len(rc.ycol))
		}
		if rc.ycol[0] >= maxRC0 {
			rc.wasimp[0] = 1.0
		}
		if rc.ycol[1] >= maxRC1 {
			rc.wasimp[1] = 1.0
		}
	}
	return rc
}

package mem

type RecoveryCycle struct{ LabTab }

var RCInterval = Column([]float64{2, 2.5, 3.2, 4, 5, 6.3, 7.9, 10, 13, 18, 24, 32, 42, 56, 75, 100, 140, 200})

func newRC() *RecoveryCycle {
	return &RecoveryCycle{LabTab{
		section:   "RECOVERY CYCLE",
		xname:     "Interval (ms)",
		yname:     "Threshold change (%)",
		xcol:      RCInterval,
		precision: 0.000001,
		logScale:  true,
	}}
}

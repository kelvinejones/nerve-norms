package mem

import (
	"math"
)

type Column []float64

func (val *Column) ImputeWithValue(oldLabel, newLabel Column, eps float64, logX bool) Column {
	num := len(newLabel)
	col := Column(make([]float64, num))
	wasImp := Column(make([]float64, num))
	colChanged := false

	intFunc := interpolate
	if logX {
		intFunc = interpolateLog
	}

	oldNum := len(*val)
	oldInd := 0
	for i, lab := range newLabel {
		for oldInd < oldNum && lab-eps > oldLabel[oldInd] {
			// The old label was for some reason not in the list of expected labels, so keep skipping until it works.
			oldInd++
		}
		if oldInd >= oldNum || lab+eps < oldLabel[oldInd] {
			// The current label is missing, so impute it with linear interpolation
			if oldNum < 2 {
				col[i] = (*val)[0]
			} else if oldInd == 0 {
				col[i] = intFunc(oldLabel[1], oldLabel[0], lab, (*val)[1], (*val)[0])
			} else if oldInd >= oldNum {
				col[i] = intFunc(oldLabel[oldNum-1], oldLabel[oldNum-2], lab, (*val)[oldNum-1], (*val)[oldNum-2])
			} else {
				col[i] = intFunc(oldLabel[oldInd], oldLabel[oldInd-1], lab, (*val)[oldInd], (*val)[oldInd-1])
			}
			wasImp[i] = 1.0
			colChanged = true
		} else {
			col[i] = (*val)[oldInd]
			oldInd++
		}
	}

	*val = col
	if colChanged {
		return wasImp
	} else {
		return Column(nil)
	}
}

func interpolate(x1, x2, x3, y1, y2 float64) float64 {
	return y2 - (x2-x3)/(x1-x2)*(y1-y2)
}

func interpolateLog(x1, x2, x3, y1, y2 float64) float64 {
	x1 = math.Log10(x1)
	x2 = math.Log10(x2)
	x3 = math.Log10(x3)
	return y2 - (x2-x3)/(x1-x2)*(y1-y2)
}

func (col Column) Maximum() float64 {
	max := math.Inf(-1)
	for _, val := range col {
		if val > max {
			max = val
		}
	}
	return max
}

func (col Column) Minimum() float64 {
	min := math.Inf(1)
	for _, val := range col {
		if val < min {
			min = val
		}
	}
	return min
}

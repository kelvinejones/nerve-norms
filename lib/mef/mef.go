package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type Mef []*mem.Mem

// Append appends the data from the second Mef, but ignores its filters. The original/combined Mef is returned.
func (mef *Mef) Append(mef2 Mef) *Mef {
	*mef = append(*mef, mef2...)
	return mef
}

func (mef *Mef) Filter(filt *Filter) *Mef {
	if filt == nil {
		// There are no filters, so return original
		return mef
	}

	mems := make(Mef, 0, len(*mef))
	for _, m := range *mef {
		// For each Mem, check if it passes all filters
		if filt.Apply(*m) {
			// This Mem passed, so append it
			mems = append(mems, m)
		}
	}

	*mef = mems
	return mef
}

package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type Mef map[string]*mem.Mem

// Append appends the data from the second Mef, but ignores its filters. The original/combined Mef is returned.
// If an exact participant name matches in both, the original is overwritten.
func (mef *Mef) Append(mef2 Mef) *Mef {
	for key, val := range mef2 {
		(*mef)[key] = val
	}
	return mef
}

func (mef *Mef) Filter(filt *Filter) *Mef {
	if filt == nil {
		// There are no filters, so return original
		return mef
	}

	mems := make(Mef)
	for key, m := range *mef {
		// For each Mem, check if it passes all filters
		if filt.Apply(*m) {
			// This Mem passed, so append it
			mems[key] = m
		}
	}

	*mef = mems
	return mef
}

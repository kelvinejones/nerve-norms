package mef

import (
	"encoding/json"

	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type Mef struct {
	mems []*mem.Mem
	Filter
}

// Append appends the data from the second Mef, but ignores its filters. The original/combined Mef is returned.
func (mef *Mef) Append(mef2 Mef) *Mef {
	mef.mems = append(mef.mems, mef2.mems...)
	return mef
}

func (mef *Mef) FilteredMef() *Mef {
	if mef.filters == nil || len(mef.filters) == 0 {
		// There are no filters, so return all data
		return mef
	}

	mems := make([]*mem.Mem, 0, len(mef.mems))
	for _, m := range mef.mems {
		// For each Mem, check if it passes all filters
		if mef.Filter.Apply(*m) {
			// This Mem passed, so append it
			mems = append(mems, m)
		}
	}

	return &Mef{mems: mems}
}

func (mef *Mef) MarshalJSON() ([]byte, error) {
	return json.Marshal(mef.FilteredMef().mems)
}

func (mef *Mef) UnmarshalJSON(value []byte) error {
	return json.Unmarshal(value, &mef.mems)
}

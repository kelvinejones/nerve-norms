package mef

import (
	"encoding/json"

	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type Mef struct {
	mems    []*mem.Mem
	filters []filter
}

// Append appends the data from the second Mef, but ignores its filters. The original/combined Mef is returned.
func (mef *Mef) Append(mef2 Mef) *Mef {
	mef.mems = append(mef.mems, mef2.mems...)
	return mef
}

func (mef *Mef) ClearFilters() {
	mef.filters = nil
}

func (mef *Mef) addFilter(filt filter) *Mef {
	mef.filters = append(mef.filters, filt)
	return mef
}

func (mef *Mef) FilterBySex(sex mem.Sex) *Mef {
	if sex == mem.UnknownSex {
		// This means no sex filtering, so don't add a filter!
		return mef
	}
	return mef.addFilter(&SexFilter{Sex: sex})
}

func (mef *Mef) FilterByAge(youngAge, oldAge int) *Mef {
	if youngAge == 0 && oldAge == 0 {
		// This means no age filtering, so don't add a filter!
		return mef
	}
	return mef.addFilter(&AgeFilter{youngAge: youngAge, oldAge: oldAge})
}

func (mef *Mef) FilteredMef() *Mef {
	if mef.filters == nil || len(mef.filters) == 0 {
		// There are no filters, so return all data
		return mef
	}

	mems := make([]*mem.Mem, 0, len(mef.mems))
	for _, m := range mef.mems {
		// For each Mem, check if it passes all filters
		shouldInclude := true
		for _, filt := range mef.filters {
			if !filt.Apply(*m) {
				// A filter was failed, so keep going.
				shouldInclude = false
				break
			}
		}
		if shouldInclude {
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

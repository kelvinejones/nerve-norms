package mef

import (
	"encoding/json"

	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type Mef struct {
	mems    []mem.Mem
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
	return mef.addFilter(&SexFilter{Sex: sex})
}

func (mef *Mef) FilterByAge(youngAge, oldAge int) *Mef {
	return mef.addFilter(&AgeFilter{youngAge: youngAge, oldAge: oldAge})
}

func (mef *Mef) MarshalJSON() ([]byte, error) {
	return json.Marshal(mef.mems)
}

func (mef *Mef) UnmarshalJSON(value []byte) error {
	return json.Unmarshal(value, &mef.mems)
}

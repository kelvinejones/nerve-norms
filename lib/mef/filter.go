package mef

import "gogs.bellstone.ca/james/jitter/lib/mem"

// filter returns true/false based on whether the provided Mem matches a filter
type filter interface {
	Filter(mem.Mem) bool
}

// filterWithConstraint updates the FilteredMef filter list based on the provided filter.
// It returns the original object to support chaining.
func (fmef *FilteredMef) filterWithConstraint(filt filter) *FilteredMef {
	for name, memData := range fmef.Mef {
		if filt.Filter(memData) {
			fmef.IncludedNames = append(fmef.IncludedNames, name)
		}
	}
	return fmef
}

// SexFilter is a type that filters sex. It uses 'mem.UnknownSex' for the unfiltered setting.
type SexFilter struct {
	mem.Sex
}

func (filt SexFilter) Filter(m mem.Mem) bool {
	return m.Header.Sex == filt.Sex || filt.Sex == mem.UnknownSex
}

// AgeFilter is a type that filters by age. It doesn't care if oldAge<youngAge, and it considers '0' to mean a value is unset.
type AgeFilter struct {
	youngAge int
	oldAge   int
}

func (filt AgeFilter) Filter(m mem.Mem) bool {
	age := m.Header.Age
	return (filt.youngAge == 0 || age >= filt.youngAge) && (filt.oldAge == 0 || age <= filt.oldAge)
}

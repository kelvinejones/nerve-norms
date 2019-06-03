package mef

import "gogs.bellstone.ca/james/jitter/lib/mem"

// filter returns true/false based on whether the provided Mem matches a filter
type filter interface {
	Apply(mem.Mem) bool
	Combine(filter) filter
}

type CompositeFilter struct {
	filters []filter
}

func (cf CompositeFilter) Combine(f filter) filter {
	cf.filters = append(cf.filters, f)
	return &cf
}

func (cf CompositeFilter) Apply(m mem.Mem) bool {
	for _, filt := range cf.filters {
		if !filt.Apply(m) {
			return false
		}
	}
	return true
}

// SexFilter is a type that filters sex. It uses 'mem.UnknownSex' for the unfiltered setting.
type SexFilter struct {
	CompositeFilter
	mem.Sex
}

func (filt SexFilter) Apply(m mem.Mem) bool {
	return m.Header.Sex == filt.Sex || filt.Sex == mem.UnknownSex
}

// AgeFilter is a type that filters by age. It doesn't care if oldAge<youngAge, and it considers '0' to mean a value is unset.
type AgeFilter struct {
	CompositeFilter
	youngAge int
	oldAge   int
}

func (filt AgeFilter) Apply(m mem.Mem) bool {
	age := m.Header.Age
	return (filt.youngAge == 0 || age >= filt.youngAge) && (filt.oldAge == 0 || age <= filt.oldAge)
}

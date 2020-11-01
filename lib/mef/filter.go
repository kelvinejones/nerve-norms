package mef

import "github.com/GrantJLiu/nerve-norms/lib/mem"

type Filter struct {
	filters []filter
}

func NewFilter() *Filter {
	return &Filter{}
}

func (cf *Filter) add(f filter) *Filter {
	cf.filters = append(cf.filters, f)
	return cf
}

func (cf Filter) Apply(m mem.Mem) bool {
	for _, filt := range cf.filters {
		if !filt.Apply(m) {
			return false
		}
	}
	return true
}

// filter returns true/false based on whether the provided Mem matches a filter
type filter interface {
	Apply(mem.Mem) bool
}

// SexFilter is a type that filters sex. It uses 'mem.UnknownSex' for the unfiltered setting.
type SexFilter struct {
	mem.Sex
}

func (filt SexFilter) Apply(m mem.Mem) bool {
	return m.Header.Sex == filt.Sex || filt.Sex == mem.UnknownSex
}

func (cf *Filter) BySex(sex mem.Sex) *Filter {
	if sex == mem.UnknownSex {
		// This means no sex filtering, so don't add a filter!
		return cf
	}
	return cf.add(&SexFilter{Sex: sex})
}

// AgeFilter is a type that filters by age. It doesn't care if oldAge<youngAge, and it considers '0' to mean a value is unset.
type AgeFilter struct {
	youngAge int
	oldAge   int
}

func (filt AgeFilter) Apply(m mem.Mem) bool {
	age := m.Header.Age
	return (filt.youngAge == 0 || age >= filt.youngAge) && (filt.oldAge == 0 || age <= filt.oldAge)
}

func (cf *Filter) ByAge(youngAge, oldAge int) *Filter {
	if youngAge == 0 && oldAge == 0 {
		// This means no age filtering, so don't add a filter!
		return cf
	}
	return cf.add(&AgeFilter{youngAge: youngAge, oldAge: oldAge})
}

// CountryFilter is a type that filters by country.
type CountryFilter struct {
	country string
}

func (filt CountryFilter) Apply(m mem.Mem) bool {
	country := m.Header.Country
	return country == filt.country || "" == filt.country
}

func (cf *Filter) ByCountry(country string) *Filter {
	if country == "" {
		// This means no country filtering, so don't add a filter!
		return cf
	}
	return cf.add(&CountryFilter{country: country})
}

// SpeciesFilter is a type that filters by species.
type SpeciesFilter struct {
	species string
}

func (filt SpeciesFilter) Apply(m mem.Mem) bool {
	species := m.Header.Species
	return species == filt.species || "" == filt.species
}

func (cf *Filter) BySpecies(species string) *Filter {
	if species == "" {
		// This means no species filtering, so don't add a filter!
		return cf
	}
	return cf.add(&SpeciesFilter{species: species})
}

// NerveFilter is a type that filters by nerve.
type NerveFilter struct {
	nerve string
}

func (filt NerveFilter) Apply(m mem.Mem) bool {
	nerve := m.Header.Nerve
	return nerve == filt.nerve || "" == filt.nerve
}

func (cf *Filter) ByNerve(nerve string) *Filter {
	if nerve == "" {
		// This means no nerve filtering, so don't add a filter!
		return cf
	}
	return cf.add(&NerveFilter{nerve: nerve})
}

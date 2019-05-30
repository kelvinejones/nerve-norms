package mef

import "gogs.bellstone.ca/james/jitter/lib/mem"

type Mef []mem.Mem

type FilteredMef struct {
	Mef
	IncludedNames []string
}

func (mef Mef) FilteredBySex(sex mem.Sex) FilteredMef {
	return *(&FilteredMef{Mef: mef}).filterWithConstraint(SexFilter{Sex: sex})
}

func (mef Mef) FilteredByAge(youngAge, oldAge int) FilteredMef {
	return *(&FilteredMef{Mef: mef}).filterWithConstraint(AgeFilter{youngAge: youngAge, oldAge: oldAge})
}

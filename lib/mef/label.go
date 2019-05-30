package mef

func (mef *Mef) LabelWithCountry(country string) *Mef {
	for _, memData := range mef.mems {
		memData.Header.Country = country
	}
	return mef
}

func (mef *Mef) LabelWithSpecies(species string) *Mef {
	for _, memData := range mef.mems {
		memData.Header.Species = species
	}
	return mef
}

func (mef *Mef) LabelWithNerve(nerve string) *Mef {
	for _, memData := range mef.mems {
		memData.Header.Nerve = nerve
	}
	return mef
}

package mef

type Norm struct {
	CDNorm GenericNorm `json:"cd"`
	RCNorm GenericNorm `json:"rc"`
	SRNorm `json:"sr"`
	IVNorm GenericNorm `json:"iv"`
	TENorm `json:"te"`
}

func (mef *Mef) Norm() Norm {
	mef = mef.FilteredMef()

	return Norm{
		CDNorm: mef.cdNorm(),
		RCNorm: mef.rcNorm(),
		SRNorm: mef.srNorm(),
		IVNorm: mef.ivNorm(),
		TENorm: mef.teNorm(),
	}
}

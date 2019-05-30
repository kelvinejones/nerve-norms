package mef

type Norm struct {
	CDNorm
	RCNorm
	SRNorm
}

func (mef *Mef) Norm() Norm {
	mef = mef.FilteredMef()

	return Norm{
		CDNorm: mef.cdNorm(),
		RCNorm: mef.rcNorm(),
		SRNorm: mef.srNorm(),
	}
}

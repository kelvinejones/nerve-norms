package mef

type Norm struct {
	CDNorm
	RCNorm
	SRNorm
	IVNorm
}

func (mef *Mef) Norm() Norm {
	mef = mef.FilteredMef()

	return Norm{
		CDNorm: mef.cdNorm(),
		RCNorm: mef.rcNorm(),
		SRNorm: mef.srNorm(),
		IVNorm: mef.ivNorm(),
	}
}

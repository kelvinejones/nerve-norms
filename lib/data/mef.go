package data

import (
	"encoding/json"

	"gogs.bellstone.ca/james/jitter/lib/mef"
)

func AsMef() (mef.Mef, error) {
	rawIn := json.RawMessage(Participants)
	bytes, err := rawIn.MarshalJSON()
	if err != nil {
		return mef.Mef{}, err
	}

	var mefData mef.Mef
	err = json.Unmarshal(bytes, &mefData)
	return mefData, nil
}

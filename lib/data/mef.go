package data

import (
	"encoding/json"

	"gogs.bellstone.ca/james/jitter/lib/mef"
)

func AsMef() (mef.Mef, error) {
	var mefData mef.Mef
	err := json.Unmarshal(jsonMef, &mefData)
	return mefData, err
}

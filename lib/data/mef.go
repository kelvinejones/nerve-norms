package data

import (
	"encoding/json"

	"github.com/GrantJLiu/nerve-norms/lib/mef"
)

func AsMef() (mef.Mef, error) {
	var mefData mef.Mef
	err := json.Unmarshal(jsonMef, &mefData)
	return mefData, err
}

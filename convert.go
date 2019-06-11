package jitter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gogs.bellstone.ca/james/jitter/lib/data"
	"gogs.bellstone.ca/james/jitter/lib/mef"
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type memWithScores struct {
	*mem.Mem      `json:"participant"`
	mef.OutScores `json:"outlierScores"`
}

func ConvertMemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	mefData, err := data.AsMef()
	if err != nil {
		setError(w, "Error loading MEF: "+err.Error())
		return
	}

	norm, err := getFilteredNormsFromRequest(r, mefData)
	if err != nil {
		setError(w, "Error getting filtered norms due to "+err.Error())
		return
	}

	if r.Body == nil {
		setError(w, "Could not load participant from empty body")
		return
	}

	var mws memWithScores
	mws.Mem, err = mem.Import(r.Body)
	if err != nil {
		setError(w, "Could not load Mem from body because "+err.Error())
		return
	}

	mws.OutScores = norm.OutlierScores(mws.Mem)

	jsOSArray, err := json.Marshal(&mws)
	if err != nil {
		setError(w, "Could not create outlier score JSON due to error: "+err.Error())
		return
	}

	log.Println("Served converted MEM for " + mws.Mem.Header.Name)
	fmt.Fprintln(w, string(jsOSArray))
}

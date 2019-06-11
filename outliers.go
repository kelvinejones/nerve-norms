package jitter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gogs.bellstone.ca/james/jitter/lib/data"
)

func OutlierScoreHandler(w http.ResponseWriter, r *http.Request) {
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

	name := r.FormValue("name")
	if name == "" {
		setError(w, "Could not make outlier scores for empty participant")
		return
	}

	mm := mefData.MemWithKey(name)
	if mm == nil {
		setError(w, "Could not find participant '"+name+"'")
		return
	}

	os := norm.OutlierScores(mm)

	jsOSArray, err := json.Marshal(&os)
	if err != nil {
		setError(w, "Could not create outlier score JSON due to error: "+err.Error())
		return
	}
	log.Println("Served outlier scores for " + name)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintln(w, string(jsOSArray))
}

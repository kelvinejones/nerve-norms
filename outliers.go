package jitter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gogs.bellstone.ca/james/jitter/lib/data"
	"gogs.bellstone.ca/james/jitter/lib/mef"
)

func OutlierScoreHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		setError(w, "Error parsing form: "+err.Error())
		return
	}
	name := r.FormValue("name")
	if name == "" {
		setError(w, "Could not make outlier scores for empty participant")
		return
	}

	mefData, err := data.AsMef()
	if err != nil {
		setError(w, "Error loading MEF: "+err.Error())
		return
	}

	mm := mefData.MemWithKey(name)
	if mm == nil {
		setError(w, "Could not find participant '"+name+"'")
		return
	}

	sex, err := parseSex(r.FormValue("sex"))
	if err != nil {
		setError(w, "Error parsing sex: "+err.Error())
		return
	}
	minAge, err := strconv.Atoi(r.FormValue("minAge"))
	if err != nil {
		setError(w, "Could not parse minAge due to error: "+err.Error())
		return
	}
	maxAge, err := strconv.Atoi(r.FormValue("maxAge"))
	if err != nil {
		setError(w, "Could not parse minAge due to error: "+err.Error())
		return
	}

	mefData.Filter(mef.NewFilter().BySex(sex).ByAge(minAge, maxAge).ByCountry(r.FormValue("country")).BySpecies(r.FormValue("species")).ByNerve(r.FormValue("nerve")))
	norm := mefData.Norm()
	os := norm.OutlierScores(mm)

	jsOSArray, err := json.Marshal(&os)
	if err != nil {
		setError(w, "Could not create outlier score JSON due to error: "+err.Error())
		return
	}
	log.Println("Served outlier scores for " + name)
	fmt.Fprintln(w, jsOSArray)
}

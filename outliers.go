package jitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/GrantJLiu/nerve-norms/lib/data"
	"github.com/GrantJLiu/nerve-norms/lib/mef"
	"github.com/GrantJLiu/nerve-norms/lib/mem"
)

func OutlierScoreHandler(w http.ResponseWriter, r *http.Request) {
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

	name, mm, err := getMemFromRequest(r, mefData)
	if err != nil {
		setError(w, "Could not load Mem from request because "+err.Error())
		return
	}

	os := norm.OutlierScores(mm)
	jsOSArray, err := json.Marshal(&os)
	if err != nil {
		setError(w, "Could not create outlier score JSON due to error: "+err.Error())
		return
	}

	log.Println("Served outlier scores for " + name)
	fmt.Fprintln(w, string(jsOSArray))
}

func getMemFromRequest(r *http.Request, mefData mef.Mef) (string, *mem.Mem, error) {
	name := r.FormValue("name")
	if name == "" {
		if r.Body == nil {
			return "", nil, errors.New("could not make outlier scores for empty participant")
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return "", nil, errors.New("could not read request body because " + err.Error())
		}

		var mm mem.Mem
		err = json.Unmarshal(body, &mm)
		if err != nil {
			return "", nil, errors.New("could not unmarshal Mem because " + err.Error())
		}

		return "Uploaded MEM '" + mm.Header.Name + "'", &mm, nil
	}

	mm := mefData.MemWithKey(name)
	if mm == nil {
		return "", nil, errors.New("could not find participant '" + name + "'")
	}

	return name, mm, nil
}

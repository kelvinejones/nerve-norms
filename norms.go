package jitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gogs.bellstone.ca/james/jitter/lib/data"
	"gogs.bellstone.ca/james/jitter/lib/mef"
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

func NormHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		setError(w, "Error parsing form: "+err.Error())
		return
	}

	var mefData mef.Mef
	err = json.Unmarshal([]byte(data.Participants), &mefData)
	if err != nil {
		setError(w, "Error unmarshaling participants: "+err.Error())
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

	jsNorm := mefData.Norm()
	jsNormArray, err := json.Marshal(&jsNorm)
	if err != nil {
		setError(w, "Could not create norm JSON due to error: "+err.Error())
		return
	}
	log.Println("Served norms")
	fmt.Fprintln(w, jsNormArray)
}

func setError(w http.ResponseWriter, str string) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println(str)
	fmt.Fprintln(w, str)
}

func parseSex(sex string) (mem.Sex, error) {
	switch sex {
	case "male", "Male", "M", "m":
		return mem.MaleSex, nil
	case "female", "Female", "F", "f":
		return mem.MaleSex, nil
	case "":
		return mem.UnknownSex, nil
	default:
		return mem.UnknownSex, errors.New("Invalid sex '" + sex + "'")
	}
}

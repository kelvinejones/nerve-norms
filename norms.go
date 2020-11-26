package jitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/GrantJLiu/nerve-norms/lib/data"
	"github.com/GrantJLiu/nerve-norms/lib/mef"
	"github.com/GrantJLiu/nerve-norms/lib/mem"
)

type filterParameters struct {
	sex     mem.Sex
	minAge  int
	maxAge  int
	country string
	species string
	nerve   string
}

// NormHandler API endpoint for requests about the normative data for applied filters
func NormHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	mefData, err := data.AsMef()
	if err != nil {
		setError(w, "Error loading MEF because "+err.Error())
		return
	}

	jsNorm, err := getFilteredNormsFromRequest(r, mefData)
	if err != nil {
		setError(w, "Error getting filtered norms due to "+err.Error())
		return
	}

	jsNormArray, err := json.Marshal(&jsNorm)
	if err != nil {
		setError(w, "Could not create norm JSON due to "+err.Error())
		return
	}
	log.Println("Served norms")
	fmt.Fprintln(w, string(jsNormArray))
}

func setError(w http.ResponseWriter, str string) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println(str)
	output := struct {
		String string `json:"error"`
	}{str}

	js, err := json.Marshal(&output)
	if err != nil {
		log.Println("Could not marshal error " + err.Error())
		fmt.Fprintln(w, str)
	} else {
		fmt.Fprintln(w, string(js))
	}
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

func parseQuery(r *http.Request) (filterParameters, error) {
	fp := filterParameters{
		country: r.FormValue("country"),
		species: "human", // for now, rat doesn't work
		nerve:   r.FormValue("nerve"),
	}
	var err error
	fp.sex, err = parseSex(r.FormValue("sex"))
	if err != nil {
		return fp, errors.New("Error parsing sex: " + err.Error())
	}
	fp.minAge, err = strconv.Atoi(r.FormValue("minAge"))
	if r.FormValue("minAge") != "" && err != nil {
		return fp, errors.New("Could not parse minAge due to error: " + err.Error())
	}
	fp.maxAge, err = strconv.Atoi(r.FormValue("maxAge"))
	if r.FormValue("maxAge") != "" && err != nil {
		return fp, errors.New("Could not parse maxAge due to error: " + err.Error())
	}

	return fp, nil
}

func getFilteredNormsFromRequest(r *http.Request, mefData mef.Mef) (mef.Norm, error) {
	err := r.ParseForm()
	if err != nil {
		return mef.Norm{}, errors.New("error parsing form because " + err.Error())
	}

	fp, err := parseQuery(r)
	if err != nil {
		return mef.Norm{}, errors.New("error parsing query because " + err.Error())
	}
	mefData.Filter(mef.NewFilter().BySex(fp.sex).ByAge(fp.minAge, fp.maxAge).ByCountry(fp.country).BySpecies(fp.species).ByNerve(fp.nerve))

	return mefData.Norm(), nil
}

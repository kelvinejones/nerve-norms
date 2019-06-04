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

type filterParameters struct {
	sex     mem.Sex
	minAge  int
	maxAge  int
	country string
	species string
	nerve   string
}

func NormHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		setError(w, "Error parsing form: "+err.Error())
		return
	}

	mefData, err := data.AsMef()
	if err != nil {
		setError(w, "Error loading MEF: "+err.Error())
		return
	}

	fp, err := parseQuery(r)
	if err != nil {
		setError(w, "Error parsing query: "+err.Error())
		return
	}
	mefData.Filter(mef.NewFilter().BySex(fp.sex).ByAge(fp.minAge, fp.maxAge).ByCountry(fp.country).BySpecies(fp.species).ByNerve(fp.nerve))

	jsNorm := mefData.Norm()
	jsNormArray, err := json.Marshal(&jsNorm)
	if err != nil {
		setError(w, "Could not create norm JSON due to error: "+err.Error())
		return
	}
	log.Println("Served norms")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintln(w, string(jsNormArray))
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

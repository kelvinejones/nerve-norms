package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"gogs.bellstone.ca/james/jitter/lib/mef"
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

var input = flag.String("input", "json/all.json", "path to the JSON that should be loaded")
var output = flag.String("output", "", "path to save the filtered JSON; otherwise, do nothing with it")
var norm = flag.String("norm", "json/norm.json", "path to save the norm JSON; otherwise, output to stdout")
var sexString = flag.String("sex", "", "only include participants of this sex (M/F)")
var minAge = flag.Int("minAge", 0, "only include participants at least this old")
var maxAge = flag.Int("maxAge", 200, "only include participants this age or younger")
var country = flag.String("country", "", "only include participants from this country (CA/PO/JP)")

func main() {
	flag.Parse()

	if *minAge > *maxAge {
		panic("minAge > maxAge is not valid")
	}

	sex, err := parseSex(*sexString)
	if err != nil {
		panic(err)
	}

	file, err := os.Open(*input)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var mefData mef.Mef

	err = json.Unmarshal(bytes, &mefData)
	if err != nil {
		panic(err)
	}

	mefData.FilterBySex(sex).FilterByAge(*minAge, *maxAge)

	jsArray, err := json.Marshal(&mefData)
	if err != nil {
		fmt.Println("Could not concatenate JSON due to error: " + err.Error())
	}

	if *output != "" {
		err = ioutil.WriteFile(*output, jsArray, 0644)
		if err != nil {
			fmt.Println("Could not save JSON due to error: " + err.Error())
		}
	}

	jsNorm := mefData.Norm()
	if *norm == "" {
		fmt.Println(jsNorm)
	} else {
		jsNormArray, err := json.Marshal(&jsNorm)
		if err != nil {
			fmt.Println("Could not create norm JSON due to error: " + err.Error())
		}
		err = ioutil.WriteFile(*norm, jsNormArray, 0644)
		if err != nil {
			fmt.Println("Could not save norm JSON due to error: " + err.Error())
		}
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

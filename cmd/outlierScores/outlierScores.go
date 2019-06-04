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

var name = flag.String("name", "CA-CR21S", "name of the participant whose outlier scores are desired")
var input = flag.String("input", "json/all.json", "path to the JSON that should be loaded")
var output = flag.String("output", "", "path to save the outlier scores JSON; otherwise, print it to stdout")
var sexString = flag.String("sex", "", "only include participants of this sex (M/F) in the norms")
var minAge = flag.Int("minAge", 0, "only include participants at least this old in the norms")
var maxAge = flag.Int("maxAge", 200, "only include participants this age or younger in the norms")
var country = flag.String("country", "", "only include participants from this country (CA/PO/JP) in the norms")
var species = flag.String("species", "human", "only include participants from this species (human/rat) in the norms")
var nerve = flag.String("nerve", "median", "only include participants from this nerve (median/CP) in the norms")

func main() {
	flag.Parse()

	if *name == "" {
		panic("A participant name is required")
	}

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

	mm := mefData.MemWithKey(*name)
	if mm == nil {
		panic("Could not find participant '" + *name + "'")
	}

	mefData.Filter(mef.NewFilter().BySex(sex).ByAge(*minAge, *maxAge).ByCountry(*country).BySpecies(*species).ByNerve(*nerve))
	norm := mefData.Norm()
	os := norm.OutlierScores(mm)

	jsOSArray, err := json.Marshal(&os)
	if err != nil {
		panic("Could not create outlier score JSON due to error: " + err.Error())
	}

	if *output == "" {
		fmt.Println(string(jsOSArray))
	} else {
		err = ioutil.WriteFile(*output, jsOSArray, 0644)
		if err != nil {
			panic("Could not save outlier score JSON due to error: " + err.Error())
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

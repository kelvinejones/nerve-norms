package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"gogs.bellstone.ca/james/jitter/lib/mef"
)

var input = flag.String("input", "res/data/participants.json", "path to the JSON that should be loaded")
var output = flag.String("output", "", "path to save the filtered JSON; otherwise, output to stdout")
var sex = flag.String("sex", "", "only include participants of this sex (male/female)")
var minAge = flag.Int("minAge", 0, "only include participants at least this old")
var maxAge = flag.Int("maxAge", 200, "only include participants this age or younger")
var country = flag.String("country", "", "only include participants from this country (CA/PO/JP)")

func main() {
	flag.Parse()

	if *minAge > *maxAge {
		panic("minAge > maxAge is not valid")
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

	// TODO filter

	jsArray, err := json.Marshal(&mefData)
	if err != nil {
		fmt.Println("Could not concatenate JSON due to error: " + err.Error())
	}

	if *output == "" {
		fmt.Printf("%v\n", string(jsArray))
	} else {
		err = ioutil.WriteFile(*output, jsArray, 0644)
		if err != nil {
			fmt.Println("Could not save JSON due to error: " + err.Error())
		}
	}
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"gogs.bellstone.ca/james/jitter/lib/mef"
)

var caPath = flag.String("caPath", "/Users/james/Documents/Education/UofA/MSc/Research/normative-data/human/CA/FESmedianAPB.MEF", "path to the CA MEF")
var jpPath = flag.String("jpPath", "", "path to the JP MEF")
var poPath = flag.String("poPath", "/Users/james/Documents/Education/UofA/MSc/Research/normative-data/human/PO/Portugal.MEF", "path to the PO MEF")
var legPath = flag.String("legPath", "/Users/james/Documents/Education/UofA/MSc/Research/normative-data/human/CA/FEScommonperonealTA.MEF", "path to the leg MEF")
var ratPath = flag.String("ratPath", "", "path to the rat MEF")
var output = flag.String("output", "json/all.json", "path to save the JSON; otherwise, output to stdout")

func main() {
	flag.Parse()

	caMef, err := mef.Import(*caPath)
	if err != nil && *caPath != "" {
		panic(err)
	}
	jpMef, err := mef.Import(*jpPath)
	if err != nil && *jpPath != "" {
		panic(err)
	}
	poMef, err := mef.Import(*poPath)
	if err != nil && *poPath != "" {
		panic(err)
	}
	legMef, err := mef.Import(*legPath)
	if err != nil && *legPath != "" {
		panic(err)
	}
	ratMef, err := mef.Import(*ratPath)
	if err != nil && *ratPath != "" {
		panic(err)
	}

	allData := caMef.Append(jpMef).Append(poMef).Append(legMef).Append(ratMef)

	jsArray, err := json.Marshal(&allData)
	if err != nil {
		panic("Could not concatenate JSON due to error: " + err.Error())
	}

	if *output == "" {
		fmt.Printf("%v\n", string(jsArray))
	} else {
		err = ioutil.WriteFile(*output, jsArray, 0644)
		if err != nil {
			panic("Could not save JSON due to error: " + err.Error())
		}
	}
}

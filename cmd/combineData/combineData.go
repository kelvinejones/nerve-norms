package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"gogs.bellstone.ca/james/jitter/lib/mef"
)

var caPath = flag.String("caPath", "res/data/CA.json", "path to the CA JSON")
var jpPath = flag.String("jpPath", "", "path to the JP JSON")
var poPath = flag.String("poPath", "res/data/PO.json", "path to the PO JSON")
var legPath = flag.String("legPath", "res/data/leg.json", "path to the leg JSON")
var ratPath = flag.String("ratPath", "", "path to the rat JSON")
var output = flag.String("output", "res/data/all.json", "path to save the filtered JSON; otherwise, output to stdout")

func main() {
	caMef, err := loadJson(*caPath)
	if err != nil {
		panic(err)
	}
	jpMef, err := loadJson(*jpPath)
	if err != nil {
		panic(err)
	}
	poMef, err := loadJson(*poPath)
	if err != nil {
		panic(err)
	}
	legMef, err := loadJson(*legPath)
	if err != nil {
		panic(err)
	}
	ratMef, err := loadJson(*ratPath)
	if err != nil {
		panic(err)
	}

	allData := caMef.Append(jpMef).Append(poMef).Append(legMef).Append(ratMef)

	jsArray, err := json.Marshal(&allData)
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

func loadJson(path string) (mef.Mef, error) {
	if path == "" {
		// No data to load, so just return empty
		return mef.Mef{}, nil
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return mef.Mef{}, err
	}

	var mefData mef.Mef
	err = json.Unmarshal(bytes, &mefData)
	return mefData, err
}

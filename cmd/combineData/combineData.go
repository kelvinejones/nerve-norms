package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gogs.bellstone.ca/james/jitter/lib/mef"
)

var input = flag.String("input", "/Users/james/Documents/Education/UofA/MSc/Research/normative-data/all.csv", "path to the CSV of MEF files and info")
var output = flag.String("output", "json/all.json", "path to save the JSON; otherwise, output to stdout")
var jsFile = flag.String("jsFile", "res/templates/data/participants.json", "path to save the participants file")
var goFile = flag.String("goFile", "lib/data/data.go", "path to save a go file with the JSON")

func main() {
	flag.Parse()

	if *input == "" {
		panic("Cannot run without an imput file")
	}

	lms, err := parseLoadableMefs(*input)
	if err != nil {
		panic(err)
	}

	allData := mef.Mef{}
	meanData := mef.Mef{}
	for _, lm := range lms {
		mefData, err := mef.Import(lm.prefix, lm.path)
		if err != nil {
			panic(err)
		}
		meanData.Add(lm.name, mefData.Mean(lm.name))
		mefData.LabelWithSpecies(lm.species).LabelWithNerve(lm.nerve).LabelWithCountry(lm.country)
		allData.Append(mefData)
	}
	meanData.LabelWithSpecies("Means").LabelWithNerve("Means").LabelWithCountry("Means") // These labels make sure this data won't match any filters
	allData.Append(meanData)

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

	if *jsFile == "" && *goFile == "" {
		return
	}

	jsUnescape := strings.Replace(string(jsArray), "\\\\", "\\", -1)

	if *jsFile != "" {
		err = ioutil.WriteFile(*jsFile, []byte("const participants = "+jsUnescape+"\n"), 0644)
		if err != nil {
			panic("Could not save jsFile due to error: " + err.Error())
		}
	}

	if *goFile != "" {
		goString := "package data\n\n" +
			"var Participants = `" + jsUnescape + "`" +
			"\nvar jsonMef = []byte(`" + string(jsArray) + "`)\n"
		err = ioutil.WriteFile(*goFile, []byte(goString), 0644)
		if err != nil {
			panic("Could not save goFile due to error: " + err.Error())
		}
	}
}

type loadableMef struct {
	path    string
	prefix  string
	country string
	species string
	nerve   string
	name    string
}

func parseLoadableMefs(path string) ([]loadableMef, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvr := csv.NewReader(file)
	dir := filepath.Dir(path)

	lms := []loadableMef{}
	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return lms, err
		}
		if len(row) != 6 {
			return nil, errors.New("Incorrect number of rows in CSV")
		}

		lms = append(lms, loadableMef{
			path:    dir + "/" + row[0],
			prefix:  row[1],
			country: row[2],
			species: row[3],
			nerve:   row[4],
			name:    row[5],
		})
	}
}

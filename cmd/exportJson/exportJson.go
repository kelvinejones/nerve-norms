package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/GrantJLiu/nerve-norms/lib/mef"
	"github.com/GrantJLiu/nerve-norms/lib/mem"
)

var input = flag.String("input", "res/data/FESB70821B.MEM", "path to the file or folder that should be imported")
var output = flag.String("output", "", "name of the file to save the JSON into; if blank, print to stdout")

func main() {
	flag.Parse()

	var js []byte
	var err error
	if filepath.Ext(*input) == ".MEM" {
		js, err = loadMemAsJson(*input)
		if err != nil {
			panic("Could not load MEM due to error: " + err.Error())
		}
	} else {
		mf, err := mef.Import("", *input)
		if err != nil {
			panic("Could not load MEF due to error: " + err.Error())
		}
		js, err = json.Marshal(&mf)
		if err != nil {
			panic("Could not marshal JSON due to error: " + err.Error())
		}
	}

	if *output == "" {
		fmt.Printf("%v\n", string(js))
	} else {
		err = ioutil.WriteFile(*output, js, 0644)
		if err != nil {
			panic("Could not save JSON due to error: " + err.Error())
		}
	}
}

func loadMemAsJson(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	memData, err := mem.Import(bufio.NewReader(file))
	if err != nil {
		return nil, err
	}

	return json.Marshal(&memData)
}

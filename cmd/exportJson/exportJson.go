package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gogs.bellstone.ca/james/jitter/lib/mem"
)

var input = flag.String("input", "res/data/FESB70821B.MEM", "path to the file or folder that should be imported")
var output = flag.String("output", "", "name of the file to save the JSON into; if blank, print to stdout")

type ExpectedFiles map[string]struct{}

func main() {
	flag.Parse()

	jsonStrings := make(map[string]json.RawMessage)
	expect := ExpectedFiles(nil)

	ext := filepath.Ext(*input)
	switch ext {
	case ".MEM":
		expect := make(map[string]struct{})
		expect[*input] = struct{}{}
	case ".MEF":
		var err error
		expect, err = loadFilenamesFromMEF(*input)
		if err != nil {
			panic("Loaded bad MEF: " + err.Error())
		}
		*input = filepath.Dir(*input)
	}

	// This works regardless of whether *input is a file or a directory, though if it's an invalid file, it will be silently skipped
	filepath.Walk(*input, func(subpath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Could not walk '" + subpath + "' due to error: " + err.Error())
			return nil
		}

		basename := filepath.Base(subpath)
		// If we're expecting specific files, make sure we only look for those ones
		if expect != nil && !expect.Contains(basename) {
			return nil // Skip it
		}

		if !strings.Contains(subpath, ".MEM") || info.Mode().IsDir() {
			// Skip directories and invalid files
			return nil
		}

		name, js, err := loadMemAsJson(subpath)
		if err != nil {
			fmt.Println("Could not parse '" + basename + "' due to error: " + err.Error())
			return nil
		}
		if _, ok := jsonStrings[name]; ok {
			fmt.Println("Warning: Participant '" + name + "' already exists and has been overwritten")
		}
		jsonStrings[name] = js

		if expect != nil {
			expect.Remove(basename)
		}

		return nil
	})

	if expect != nil && len(expect) > 0 {
		fmt.Println("Not all expected files were found", expect)
	}

	jsArray, err := json.Marshal(&jsonStrings)
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

func loadMemAsJson(path string) (string, []byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", nil, err
	}

	memData, err := mem.Import(bufio.NewReader(file))
	if err != nil {
		return "", nil, err
	}

	err = memData.Verify()
	if err != nil {
		return "", nil, err
	}

	js, err := json.Marshal(&memData)
	return memData.Header.Name, js, err
}

func (ef ExpectedFiles) Contains(path string) bool {
	_, ok := ef[path]
	return ok
}

func (ef *ExpectedFiles) Remove(path string) {
	delete(*ef, path)
}

func loadFilenamesFromMEF(path string) (ExpectedFiles, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	expect := make(map[string]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		expect[scanner.Text()+".MEM"] = struct{}{}
	}

	return expect, scanner.Err()
}

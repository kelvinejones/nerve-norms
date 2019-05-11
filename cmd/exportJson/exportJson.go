package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gogs.bellstone.ca/james/jitter/lib/mem"
)

var path = flag.String("path", "res/data/FESB70821B.MEM", "path to the file that should be imported")

func main() {
	flag.Parse()

	jsonStrings := make(map[string]json.RawMessage)

	// This works regardless of whether *path is a file or a directory, though if it's an invalid file, it will be silently skipped
	filepath.Walk(*path, func(subpath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Could not walk '" + subpath + "' due to error: " + err.Error())
			return nil
		}

		if !strings.Contains(subpath, ".MEM") || info.Mode().IsDir() {
			// Skip directories and invalid files
			return nil
		}

		name, js, err := printMem(subpath)
		if err != nil {
			fmt.Println("Could not parse '" + subpath + "' due to error: " + err.Error())
		}
		jsonStrings[name] = js

		return nil
	})

	jsArray, err := json.Marshal(&jsonStrings)
	if err != nil {
		fmt.Println("Could not concatenate JSON due to error: " + err.Error())
	}

	fmt.Printf("%v\n", string(jsArray))
}

func printMem(path string) (string, []byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", nil, err
	}

	memData, err := mem.Import(bufio.NewReader(file))
	if err != nil {
		return "", nil, err
	}

	js, err := json.Marshal(&memData)
	return memData.Header.Name, js, err
}

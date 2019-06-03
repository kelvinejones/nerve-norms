package mef

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type ExpectedFiles map[string]struct{}

// Import imports the MEF at the provided path, assuming all MEM are in the same directory.
// If the input path is not an MEF, then assume it's a directory and import all
func Import(input string) (Mef, error) {
	var expect ExpectedFiles
	fi, err := os.Stat(input)
	switch {
	case filepath.Ext(input) == ".MEF":
		expect, err = loadFilenamesFromMEF(input)
		if err != nil {
			return Mef{}, errors.New("Loaded bad MEF: " + err.Error())
		}
		input = filepath.Dir(input)
	case err == nil && fi.IsDir():
		// Do nothing
	default:
		return Mef{}, errors.New("Provided import path is neither MEF nor directory")
	}

	mf := Mef{}

	// This works regardless of whether input is a file or a directory, though if it's an invalid file, it will be silently skipped
	filepath.Walk(input, func(subpath string, info os.FileInfo, err error) error {
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

		mm, err := loadMem(subpath)
		if err != nil {
			fmt.Println("Could not parse '" + basename + "' due to error: " + err.Error())
			return nil
		}
		mf = append(mf, mm)

		if expect != nil {
			expect.Remove(basename)
		}

		return nil
	})

	if expect != nil && len(expect) > 0 {
		return mf, fmt.Errorf("Not all expected files were found", expect)
	}

	return mf, nil
}

func loadMem(path string) (*mem.Mem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return mem.Import(bufio.NewReader(file))
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

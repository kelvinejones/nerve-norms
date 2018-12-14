package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gogs.bellstone.ca/james/jitter/lib/mem"
)

var path = flag.String("path", "res/data/short_test.MEM", "path to the file or directory that should be imported")

func main() {
	flag.Parse()

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

		if err := printMem(subpath); err != nil {
			fmt.Println("Could not parse '" + subpath + "' due to error: " + err.Error())
		}
		return nil
	})
}

func printMem(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	memData, err := mem.Import(bufio.NewReader(file))
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", memData)

	return nil
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"gogs.bellstone.ca/james/jitter/lib/mem"
)

var path = flag.String("path", "res/data/short_test.MEM", "path to the file that should be imported")

func main() {
	flag.Parse()

	err := printMem(*path)
	if err != nil {
		fmt.Println("Could not parse '" + *path + "' due to error: " + err.Error())
	}
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

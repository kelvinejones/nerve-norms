package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"gogs.bellstone.ca/james/jitter/lib/mem"
)

var path = flag.String("path", "res/data/FESB70821B.MEM", "path to the file that should be imported")

func main() {
	flag.Parse()

	if !strings.Contains(*path, ".MEM") {
		panic("Invalid path '" + *path + "'")
	}

	file, err := os.Open(*path)
	if err != nil {
		fmt.Println("Could not open '" + *path + "' due to error: " + err.Error())
	}

	memData, err := mem.Import(bufio.NewReader(file))
	if err != nil {
		fmt.Println("Could not parse '" + *path + "' due to error: " + err.Error())
	}

	js, err := json.Marshal(&memData)
	if err != nil {
		fmt.Println("Could not marshal JSON at '" + *path + "' due to error: " + err.Error())
	}
	fmt.Printf("%v\n", string(js))
}

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

	file, err := os.Open(*path)
	if err != nil {
		panic(err)
	}

	memData, err := mem.Import(bufio.NewReader(file))
	if err != nil {
		panic(err)
	}
	fmt.Println(memData)
}

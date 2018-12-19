package mem

/*

The mem package is used to load and store MEM files.
* `mem.go` contains the main type for MEM files and code to import them. `mem.Import` is the main function to be used from this package.
* `section.go` contains code to import sections (other than header and variables).
* `header.go` imports an MEM header.
* `variables.go` imports variables from the MEM.
* `reader.go` reads a MEM file line by line.
* `types.go` contains miscellaneous types.
* `stimulus_response.go` contains API for convenient access to stimulus-response results.

This package currently exports a lot of things that shouldn't be exported.

*/

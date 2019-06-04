package version

import (
	"fmt"
	"log"
	"net/http"
)

const version = 0.1

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Served version number", version)
	fmt.Fprintln(w, fmt.Sprintf("Version %f", version))
}

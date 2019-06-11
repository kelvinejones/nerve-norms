package jitter

import (
	"fmt"
	"log"
	"net/http"

	"gogs.bellstone.ca/james/jitter/lib/data"
)

func ParticipantHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	log.Println("Served participants")
	fmt.Fprintln(w, data.Participants)
}

package jitter

import (
	"fmt"
	"log"
	"net/http"

	"gogs.bellstone.ca/james/jitter/lib/data"
)

func ParticipantHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Served participants")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintln(w, data.Participants)
}

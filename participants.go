package jitter

import (
	"fmt"
	"log"
	"net/http"

	"gogs.bellstone.ca/james/jitter/lib/data"
)

func ParticipantHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Served participants")
	fmt.Fprintln(w, data.Participants)
}

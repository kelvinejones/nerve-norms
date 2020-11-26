package jitter

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GrantJLiu/nerve-norms/lib/data"
)

// ParticipantHandler Not sure what it does, depreciated? Participants are served to client directly in a json
func ParticipantHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	log.Println("Served participants")
	fmt.Fprintln(w, data.Participants)
}

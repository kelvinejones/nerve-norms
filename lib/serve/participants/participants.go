package participants

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func ParticipantHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("../../../res/templates/data/participants.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		fmt.Fprintln(w, err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		fmt.Fprintln(w, err)
	}

	log.Println("Served participants")
	fmt.Fprintln(w, bytes)
}

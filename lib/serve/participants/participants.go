package participants

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func ParticipantHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("res/templates/data/participants.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(w, bytes)
}

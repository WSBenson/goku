package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type fighter struct {
	Name string `json:"name"`
}

// struct to hold the new concatination of the user's name
type allegiance struct {
	Message string `json:"message"`
}

// gokuHandler accepts a post request that takes the user's name from the JSON request
// and responds back with that user's name + some message
func handleGokuRequests(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleGokuPost(w, r)

	case http.MethodGet:
		handleGokuGet(w, r)

	// all other cases
	default:
		handleGokuDefault(w, r)
	}

}

func handleGokuPost(w http.ResponseWriter, r *http.Request) {
	// Reads the body of the JSON request into the fighter struct
	// by unmarshaling the JSON bytes
	f := fighter{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &f)

	// Concatinates the name from the POST body with the string
	// then marshals the response back into JSON bytes to
	// send to the browser.
	a := allegiance{
		Message: f.Name + " has joined the Z squad",
	}
	d, err := json.Marshal(a)
	if err != nil {
		log.Fatal(err)
	}

	// Writes the JSON bytes to the JSON response writer
	w.Write(d)
}

func handleGokuGet(w http.ResponseWriter, r *http.Request) {
	// writes "The Z fighters are ready to fight"
	a := allegiance{
		Message: "The Z fighters are ready to battle",
	}
	d, err := json.Marshal(a)
	if err != nil {
		log.Fatal(err)
	}

	// Writes the JSON bytes to the JSON response writer
	w.Write(d)
}

func handleGokuDefault(w http.ResponseWriter, r *http.Request) {
	// writes "The Z fighters perished"
	a := allegiance{
		Message: "The Z fighters perished",
	}
	d, err := json.Marshal(a)
	if err != nil {
		log.Fatal(err)
	}

	// Writes the JSON bytes to the JSON response writer
	w.Write(d)
}

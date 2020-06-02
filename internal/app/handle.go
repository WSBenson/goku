package app

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
	logger = zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339},
	).With().Timestamp().Caller().Logger()
)

type fighters struct {
	Fighters []fighter `json:"fighters"`
}

type fighter struct {
	Name  string `json:"name"`
	Power int    `json:"power"`
}

// struct to hold the new concatination of the user's name
type allegiance struct {
	Message string `json:"message"`
}

// handleGokuRequests accepts a post request that takes the user's name from the JSON request
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
	f := fighters{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &f)

	// Concatinates the name from the POST body with a string
	// then marshals the response back into JSON bytes to
	// send to the browser.
	a := allegiance{}
	switch length := len(f.Fighters); {
	case length == 1:
		lenOfOne(f.Fighters[0], &a)

	case length == 2:
		lenOfTwo(f.Fighters[0], f.Fighters[1], &a)

	case length > 2:
		lenOverTwo(f, &a)

	default:
		logger.Error().Msg("Improperly formated JSON message")
		a.Message = "You did not format your JSON message properly, check the console log and the README."
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

// lenOfOne handles what JSON message to give back for a JSON array
// with one fighter in it
func lenOfOne(f fighter, a *allegiance) {
	if f.Name == "" {
		logger.Error().Msg("Some fighter name JSON fields not specified, or empty.")
		a.Message = "Some fighter name JSON fields not specified, or empty. Check the console log and the README."
		return
	}
	if f.Power == 0 {
		logger.Error().Msg("Fighter power JSON field not specified, or 0.")
		a.Message = "Fighter power JSON field not specified, or 0. Check the console log and the README."
		return
	} else if f.Power > 9000 {
		a.Message = "The scouter says " + f.Name + "'s power level is over 9000! You better start running."
	} else if f.Power < 206 {
		a.Message = f.Name + " is weaker than Krillin, cmon bruh."
	} else {
		a.Message = f.Name + "'s power level isn't even over 9000, they're a straight side character."
	}
}

// lenOfTwo handles what JSON message to give back for a JSON array
// with two fighters in it
func lenOfTwo(f fighter, f1 fighter, a *allegiance) {
	if f.Name == "" || f1.Name == "" {
		logger.Error().Msg("Some fighter name JSON fields not specified, or empty.")
		a.Message = "Some fighter name JSON fields not specified, or empty. Check the console log and the README."
		return
	}
	if f.Power == 0 || f1.Power == 0 {
		logger.Error().Msg("Some fighter power JSON fields not specified, or 0.")
		a.Message = "Some fighter power JSON fields not specified, or 0. Check the console log and the README."
		return
	}
	if f.Power == f1.Power && f.Power < 206 && f1.Power < 206 {
		a.Message = f.Name + " and " + f1.Name + " are equally trash, they better fuse or something."
	} else if f.Power < 206 && f1.Power < 206 {
		a.Message = f.Name + " and " + f1.Name + " are both weaker than Krillin, im done."
	} else if f.Power > f1.Power {
		a.Message = f.Name + "'s power level is superior to " + f1.Name + "'s"
	} else if f.Power == f1.Power {
		a.Message = f.Name + "'s power level is equal to " + f1.Name + "'s"
	} else {
		a.Message = f1.Name + "'s power level is superior to " + f.Name + "'s"
	}
}

//lenOverTwo handles what JSON message to give back for a JSON array
// with more than two fighters in it
func lenOverTwo(fs fighters, a *allegiance) {
	maxPower := 0
	maxFighter := fs.Fighters[0]
	// loops through each fighter to compare their powers
	for _, fighter := range fs.Fighters {
		if fighter.Name == "" {
			logger.Error().Msg("Some fighter name JSON fields not specified, or empty.")
			a.Message = "Some fighter name JSON fields not specified, or empty. Check the console log and the README."
			return
		}
		if fighter.Power == 0 {
			logger.Error().Msg("Some fighter power JSON fields not specified, or 0.")
			a.Message = "Some fighter power JSON fields not specified, or 0. Check the console log and the README."
			return
		}
		if fighter.Power > maxPower {
			// Sets the highest power and the fighter who has the highest power
			maxPower = fighter.Power
			maxFighter = fighter
		}
	}

	a.Message = maxFighter.Name + " is the strongest of all the Z fighters."
}

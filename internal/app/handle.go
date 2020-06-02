package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
	logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Caller().Logger()
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
		logger.Info().Msg("Goku POST request received")
		handleGokuPost(w, r)

	case http.MethodGet:
		logger.Info().Msg("Goku GET request received")
		handleGokuGet(w, r)

	// all other cases
	default:
		logger.Info().Msg("Goku other request received")
		handleGokuDefault(w, r)
	}

}

func handleGokuPost(w http.ResponseWriter, r *http.Request) {
	// Reads the body of the JSON request into the fighter struct
	// by unmarshaling the JSON bytes
	f := fighters{}
	logger.Info().Msg("Reading POST request...")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to read POST request body.")
	}
	logger.Info().Msg("Unmarshaling POST request...")
	err = json.Unmarshal(body, &f)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to unmarshal POST request body.")
	}

	// The struct that will hold the JSON message that is sent back to the browser
	a := allegiance{}

	// Checks for any JSON name and power fields that are empty, or zero.
	// (Not that important so this could be removed, I just prefer names to not be "" and powers to not be 0)
	for _, fighter := range f.Fighters {
		if fighter.Name == "" {
			logger.Error().Msg("Some fighter name JSON fields are empty.\n")
			a.Message = "Some fighter name JSON fields are empty. Check the console log and the README."
			return
		}
		if fighter.Power == 0 {
			logger.Error().Msg("Some fighter power JSON fields are 0.\n")
			a.Message = "Some fighter power JSON fields are 0. Check the console log and the README."
			return
		}
	}

	// The gokuPOSTCases function evaluates how many fighters are in the JSON POST body
	// and returns a concatinated string with the name of the fighter(s) and other information
	a.Message = gokuPOSTCases(f)

	// The response is marshaled back into JSON bytes to be sent to the browser
	logger.Info().Msg("Marshaling POST message...")
	d, err := json.Marshal(a)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to marshal POST message.")
	}

	// Writes the JSON bytes to the JSON response writer
	w.Write(d)
	logger.Info().Msg("POST message written\n")
}

func handleGokuGet(w http.ResponseWriter, r *http.Request) {
	// writes "The Z fighters are ready to fight"
	a := allegiance{
		Message: "The Z fighters are ready to battle",
	}
	logger.Info().Msg("Marshaling GET message...")
	d, err := json.Marshal(a)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to marshal GET message.")
	}

	// Writes the JSON bytes to the JSON response writer
	w.Write(d)
	logger.Info().Msg("GET message written\n")
}

func handleGokuDefault(w http.ResponseWriter, r *http.Request) {
	// writes "The Z fighters perished"
	a := allegiance{
		Message: "The Z fighters perished",
	}
	logger.Info().Msg("Marshaling message...")
	d, err := json.Marshal(a)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to marshal message.")
	}

	// Writes the JSON bytes to the JSON response writer
	w.Write(d)
	logger.Info().Msg("Message written\n")
}

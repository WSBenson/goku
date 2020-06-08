package app

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/olivere/elastic"
	"github.com/rs/zerolog"
)

// This is the mapping I chose...
// Defines which fields are stored and indexed within the elasticsearch fighters index.
// Under mappings, the fighter line designates the type of the document.
// Each section under properties are fields that the document depends on.
// The keyword field for the name property will be great for searching, sorting and
// grouping names.
const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"fighter":{
			"properties":{
				"name":{
					"type":"keyword"
				},
				"power":{
					"type":"integer"
				},
				"suggest_field":{
					"type":"completion"
				}
			}
		}
	}
}`

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

	//ctx, client :=
	// Calls the elasticClient function to create an elastic search client and
	// add an index named fighters that will use the mapping variable to set
	// the layout of its body.
	elasticClient()

	// Will be used later to add different fighters to the elastic search index
	// for _, fighter := range f.Fighters {
	// 	elasAddSearch(ctx, client, fighter)
	// }

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

// the elasticClient function creates an elastic search client and
// adds an index named fighters that will use the mapping variable to set
// the layout of its body if it doesn't already exist.
func elasticClient() (context.Context, *elastic.Client) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Caller().Timestamp().Logger()

	// Passing a context to execute each service
	ctx := context.Background()

	// Obtain a client and connect to the default Elasticsearch installation
	// on es:9200.
	client, err := elastic.NewClient(elastic.SetURL("http://es:9200"))
	if err != nil {
		// Handle error
		logger.Fatal().Err(err).Msg("failed to make new elastic search client")
	}

	// Use the IndexExists service to check if the specified fighter index exists before adding it.
	exists, err := client.IndexExists("fighters").Do(ctx)
	if err != nil {
		// Handle error
		logger.Error().Err(err).Msg("failed to check if fighters index exists")
	}
	// If that index doesn't already exist
	if !exists {
		// Create that fighters index using the mapping variable to specify the layout of the index.
		createIndex, err := client.CreateIndex("fighters").BodyString(mapping).Do(ctx)
		if err != nil {
			// Handle error
			logger.Fatal().Err(err).Msg("failed to create new elastic search index")
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
	return ctx, client
}

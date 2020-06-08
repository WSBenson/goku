package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/olivere/elastic"
	"github.com/rs/zerolog"
)

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

	ctx, client := elasNewClient()
	for _, fighter := range f.Fighters {
		elasAddSearch(ctx, client, fighter)
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

func elasNewClient() (context.Context, *elastic.Client) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Caller().Timestamp().Logger()

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	// Obtain a client and connect to the default Elasticsearch installation
	// on es:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	client, err := elastic.NewClient(elastic.SetURL("http://es:9200"))
	if err != nil {
		// Handle error
		logger.Fatal().Err(err).Msg("failed to make new elastic search client")
	}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("fighters").Do(ctx)
	if err != nil {
		// Handle error
		logger.Error().Err(err).Msg("failed to check if fighters index exists")
	}
	if !exists {
		// Create a new index.
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

func elasAddSearch(ctx context.Context, client *elastic.Client, fighter1 fighter) {
	// Index a fighter (using JSON serialization)
	put1, err := client.Index().Index("fighters").Type("fighter").Id(hashFighter(fighter1)).BodyJson(fighter1).Do(ctx)
	if err != nil {
		// Handle error
		logger.Fatal().Msg(err.Error())
	}
	logger.Info().Msgf("Indexed figher %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	// Get fighter with specified ID
	get1, err := client.Get().
		Index("fighters").
		Type("fighter").
		Id(hashFighter(fighter1)).
		Do(ctx)
	if err != nil {
		// Handle error
		logger.Fatal().Err(err).Msg("law gen ki de")
	}
	if get1.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
	}

	// Flush to make sure the documents got written.
	_, err = client.Flush().Index("fighters").Do(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to flush to elastic search")
	}

	// Search with a term query
	termQuery := elastic.NewTermQuery("name", fighter1.Name)
	searchResult, err := client.Search().
		Index("fighters").  // search in index "fighters"
		Query(termQuery).   // specify the query
		Sort("name", true). // sort by "user" field, ascending
		From(0).Size(10).   // take documents 0-9
		Pretty(true).       // pretty print request and response JSON
		Do(ctx)             // execute
	if err != nil {
		// Handle error
		logger.Fatal().Err(err).Msg("failed to search")
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	var ftyp fighter
	for _, item := range searchResult.Each(reflect.TypeOf(ftyp)) {
		if t, ok := item.(fighter); ok {
			fmt.Printf("Fighter: %s, with a power level of %d\n", t.Name, t.Power)
		}
	}
	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d fighters\n", searchResult.TotalHits())
}

// hashFighter is a function for hashing the name of a fighter with sha256
// This hashed name is used for the fighter's Id.
func hashFighter(f fighter) string {
	b, _ := json.Marshal(f)

	hasher := sha256.New()
	hasher.Write(b)

	return hex.EncodeToString(hasher.Sum(nil))
}

package es

import (
	"context"
	"io/ioutil"

	"github.com/olivere/elastic"
	"github.com/spf13/viper"
)

//client :=
// Calls the elasticClient function to create an elastic search client and
// add an index named fighters that will use the mapping variable to set
// the layout of its body.
//elasticClient(ctx)

// Will be used later to add different fighters to the elastic search index
// for _, fighter := range f.Fighters {
// 	elasAddSearch(ctx, client, fighter)
// }

// the ElasticClient function creates an elastic search client and
// adds an index named fighters that will use the mapping variable to set
// the layout of its body if it doesn't already exist.
func ElasticClient(address string, ctx context.Context) *elastic.Client {
	// Reads from the mapping.json file to get the mapping variable
	b, err := ioutil.ReadFile(viper.GetString("es_mapping_file"))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to retrieve es mapping")
	}

	// This is the mapping I chose...
	// Defines which fields are stored and indexed within the elasticsearch fighters index.
	// Under mappings, the fighter line designates the type of the document.
	// Each section under properties are fields that the document depends on.
	// The keyword field for the name property will be great for searching, sorting and
	// grouping names.
	mapping := string(b)

	// Obtain a client and connect to the default Elasticsearch installation
	// on localhost:9200.
	client, err := elastic.NewClient(elastic.SetURL(address))
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
	return client
}

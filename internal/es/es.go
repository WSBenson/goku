package es

import (
	"context"
	"io/ioutil"

	"github.com/WSBenson/goku/internal"
	"github.com/olivere/elastic/v7"
)

type Client struct {
	index string
	ctx   *context.Context
	*elastic.Client
}

// NewClient creates an elastic search client and
// adds an index named fighters that will use the mapping variable to set
// the layout of its body if it doesn't already exist.
func NewClient(address, index, mappingPath string) *Client {
	internal.Logger.Debug().Msgf(address, mappingPath)

	// Reads from the mapping.json file to get the mapping variable
	mappingData, err := ioutil.ReadFile(mappingPath)
	if err != nil {
		internal.Logger.Fatal().Err(err).Msg("failed to retrieve es mapping")
	}
	mapping := string(mappingData)

	// Obtain a client and connect to the default Elasticsearch installation
	// on localhost:9200.
	client, err := elastic.NewSimpleClient(elastic.SetURL(address))
	if err != nil {
		// Handle error
		internal.Logger.Fatal().Err(err).Msg("failed to make new elastic search client")
	}
	ctx := context.Background()

	// Use the IndexExists service to check if the specified fighter index exists before adding it.
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		// Handle error
		internal.Logger.Error().Err(err).Msg("failed to check if fighters index exists")
	}
	// If that index doesn't already exist
	if !exists {
		// Create that fighters index using the mapping variable to specify the layout of the index.
		createIndex, err := client.CreateIndex(index).BodyJson(mapping).Do(ctx)
		if err != nil {
			// Handle error
			internal.Logger.Fatal().Err(err).Msg("failed to create new elastic search index")
		}
		if !createIndex.Acknowledged {
			internal.Logger.Error().Msg("elasticsearch failed to acknowledge index creation")
			// Not acknowledged
		}
		internal.Logger.Info().Msg("successfully created elasticsearch index")
	}

	return &Client{index, &ctx, client}
}

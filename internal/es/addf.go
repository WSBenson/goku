package es

import (
	"github.com/WSBenson/goku/internal"
	"github.com/WSBenson/goku/internal/fight"
)

// AddFighter adds a fighter to the elasticsearch index
// c := es.NewClient(address, mappingPath, index)
// err := c.AddFighter(f)
func (c *Client) AddFighter(f fight.Fighter) (err error) {

	internal.Logger.Debug().Msg(c.index)
	// Use the IndexExists service to check if the specified fighter index exists before adding it.
	exists, err := c.IndexExists(c.index).Do(*c.ctx)
	if err != nil {
		// Handle error
		internal.Logger.Fatal().Err(err).Msg("failed to check if fighters index exists")
	}
	if !exists {
		internal.Logger.Error().Msg("failed to find fighters client")
		return
	}

	// Index a fighter (using JSON serialization)
	_, err = c.Index().Index(c.index).Id(f.Name).BodyJson(f).Do(*c.ctx)
	if err != nil {
		internal.Logger.Fatal().Err(err).Msg("error indexing %+v fighter with %+v client")
		return
	}

	internal.Logger.Info().Msgf("indexed fighter: %+v\n", f)
	return
}

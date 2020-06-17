package es

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/WSBenson/goku/internal"
	"github.com/WSBenson/goku/internal/fight"
	"github.com/olivere/elastic/v7"
)

// QueryFighter ...
func (c *Client) QueryFighter(f fight.Fighter) (err error) {

	internal.Logger.Error().Msg(c.index)
	// Get fighters with specified document ID
	getBattle, err := c.Get().
		Index(c.index).
		Id("Battle").
		Do(*c.ctx)
	if err != nil {
		// Handle error
		internal.Logger.Fatal().Err(err).Msg("law gen ki de")
	}
	if getBattle.Found {
		internal.Logger.Info().Msgf("Got document %s in version %d from index %s, type %s\n", getBattle.Id, getBattle.Version, getBattle.Index, getBattle.Type)
	}

	// // Flush to make sure the documents got written.
	// _, err = c.Flush().Index(c.index).Do(*c.ctx)
	// if err != nil {
	// 	internal.Logger.Fatal().Err(err).Msg("failed to flush to elastic search")
	// }

	// Search with a term query
	termQuery := elastic.NewTermQuery("name", f.Name)
	searchResult, err := c.Search().
		Index(c.index).     // search in index "fighters"
		Query(termQuery).   // specify the query
		Sort("name", true). // sort by "power" field, ascending
		From(0).Size(10).   // take documents 0-9
		Pretty(true).       // pretty print request and response JSON
		Do(*c.ctx)          // execute
	if err != nil {
		// Handle error
		internal.Logger.Fatal().Err(err).Msg("failed to search")
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	var ftyp fight.Fighter
	internal.Logger.Info().Msgf("Got hits %v", searchResult.Hits.Hits)
	for _, item := range searchResult.Each(reflect.TypeOf(ftyp)) {
		if t, ok := item.(fight.Fighter); ok {
			internal.Logger.Info().Msgf("Fighter: %s, with a power level of %d\n", t.Name, t.Power)
		}
	}
	// TotalHits is another convenience function that works even when something goes wrong.
	internal.Logger.Info().Msgf("Found a total of %d fighters\n", searchResult.TotalHits())
	return
}

func (c *Client) GetFighters() (fighters fight.Fighters, err error) {

	// Do a search
	results, err := c.Search().Index(c.index).Query(elastic.NewMatchAllQuery()).Do(context.Background())
	if err != nil {
		return
	}

	for _, hit := range results.Hits.Hits {

		var f fight.Fighter
		err = json.Unmarshal(hit.Source, &f)
		if err != nil {
			return
		}
		internal.Logger.Debug().Msgf("%+v", f)
	}

	return
}

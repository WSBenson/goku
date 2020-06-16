package es

import (
	"github.com/WSBenson/goku/internal"
	"github.com/WSBenson/goku/internal/app"
)

// AddFighter adds a fighter to the elasticsearch index
// c := es.NewClient(address, mappingPath, index)
// err := c.AddFighter(f)
func (c *Client) AddFighter(f app.Fighter) (err error) {

	// Index a fighter (using JSON serialization)
	_, err = c.Index().Index(c.index).BodyJson(f).Do(*c.ctx)
	if err != nil {
		internal.Logger.Fatal().Err(err).Msg("error indexing %+v fighter with %+v client")
		return
	}

	internal.Logger.Info().Msgf("indexed fighter: %+v\n", f)
	return
}

// // hashFighter is a function for hashing the name of a fighter with sha256
// // This hashed name is used for the fighter's Id.
// func hashFighter(f app.Fighter) string {
// 	b, err := json.Marshal(f)
// 	if err != nil {
// 		internal.Logger.Fatal().Err(err).Msg("failed to marshal for hashing")
// 	}

// 	hasher := sha256.New()
// 	hasher.Write(b)

// 	return hex.EncodeToString(hasher.Sum(nil))
// }

// 	// Get fighter with specified ID
// 	get1, err := client.Get().
// 		Index("fighters").
// 		Id(hashFighter(f)).
// 		Do(ctx)
// 	if err != nil {
// 		// Handle error
// 		internal.Logger.Fatal().Err(err).Msg("law gen ki de")
// 	}
// 	if get1.Found {
// 		internal.Logger.Info().Msgf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
// 	}

// 	// Flush to make sure the documents got written.
// 	_, err = client.Flush().Index("fighters").Do(ctx)
// 	if err != nil {
// 		internal.Logger.Fatal().Err(err).Msg("failed to flush to elastic search")
// 	}

// 	// Search with a term query
// 	termQuery := elastic.NewTermQuery("name", f.Name)
// 	searchResult, err := client.Search().
// 		Index("fighters").  // search in index "fighters"
// 		Query(termQuery).   // specify the query
// 		Sort("name", true). // sort by "user" field, ascending
// 		From(0).Size(10).   // take documents 0-9
// 		Pretty(true).       // pretty print request and response JSON
// 		Do(ctx)             // execute
// 	if err != nil {
// 		// Handle error
// 		internal.Logger.Fatal().Err(err).Msg("failed to search")
// 	}

// 	// searchResult is of type SearchResult and returns hits, suggestions,
// 	// and all kinds of other information from Elasticsearch.
// 	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

// 	// Each is a convenience function that iterates over hits in a search result.
// 	// It makes sure you don't need to check for nil values in the response.
// 	// However, it ignores errors in serialization. If you want full control
// 	// over iterating the hits, see below.
// 	var ftyp app.Fighter
// 	internal.Logger.Info().Msgf("Got hits %v", searchResult.Hits.Hits)
// 	for _, item := range searchResult.Each(reflect.TypeOf(ftyp)) {
// 		if t, ok := item.(app.Fighter); ok {
// 			internal.Logger.Info().Msgf("Fighter: %s, with a power level of %d\n", t.Name, t.Power)
// 		}
// 	}
// 	// TotalHits is another convenience function that works even when something goes wrong.
// 	internal.Logger.Info().Msgf("Found a total of %d fighters\n", searchResult.TotalHits())

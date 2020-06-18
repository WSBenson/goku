package es

import (
	"encoding/json"

	"github.com/WSBenson/goku/internal"
	"github.com/WSBenson/goku/internal/fight"
	"github.com/olivere/elastic/v7"
)

// GetFighters ... queries the fighters index to unmarshal and store
// every fighter from the elastic search index to an array. Then passes
// a struct containing the array to a function that will print out a
// battle message.
func (c *Client) GetFighters() (fightersStruct fight.Fighters, err error) {
	// Query all fighters that have been added by the addf command
	results, err := c.Search().Index(c.index).Query(elastic.NewMatchAllQuery()).Do(*c.ctx)
	if err != nil {
		return
	}

	// Go through each result from the query
	for _, hit := range results.Hits.Hits {

		var f fight.Fighter
		err = json.Unmarshal(hit.Source, &f)
		if err != nil {
			return
		}
		// fighterStruct.Fighters is the struct from fight.fighters so
		// fighterStruct.Fighters is the array inside that struct
		fightersStruct.Fighters = append(fightersStruct.Fighters, f)
		internal.Logger.Debug().Msgf("%+v", f)
	}
	// calls the MessageCases function from the fight package which handles
	// what battle/fight message to display based off the struct you give it
	battleMessage := fight.MessageCases(fightersStruct)
	internal.Logger.Info().Msgf("%+v", battleMessage)

	return
}

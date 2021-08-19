package elastic

import "github.com/olivere/elastic/v7"

const GamesIndex = "games"

func NewElasticClient(esUrl string) (*elastic.Client, error) {
	cl, err := elastic.NewClient(
		elastic.SetURL(esUrl),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		return nil, err
	}

	cl.CreateIndex(GamesIndex)

	return cl, err
}
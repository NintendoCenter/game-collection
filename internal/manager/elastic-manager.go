package manager

import (
	"context"
	"encoding/json"
	"fmt"

	"NintendoCenter/game-collection/internal/protos"
	elastic2 "NintendoCenter/game-collection/internal/providers/elastic"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

type ElasticManager struct {
	client *elastic.Client
	logger *zap.Logger
}

func NewElasticManager(client *elastic.Client, logger *zap.Logger) *ElasticManager {
	return &ElasticManager{
		client: client,
		logger: logger,
	}
}

func (m *ElasticManager) IndexGame(ctx context.Context, game *protos.Game) error {
	jsonData, err := json.Marshal(game)
	if err != nil {
		m.logger.Error("cannot marshal game data")
		return err
	}

	gameJsonString := string(jsonData)
	_, err = m.client.Index().
		Index(elastic2.GamesIndex).
		BodyJson(gameJsonString).
		Do(ctx)

	return err
}

func (m *ElasticManager) SearchGames(ctx context.Context, filter *protos.FindGameRequest) ([]*protos.Game, error) {
	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("title", filter.Title))

	m.logger.Info(fmt.Sprintf("searching for '%s' via elastic", filter.Title))

	queryStr, err1 := searchSource.Source()
	queryJs, err2 := json.Marshal(queryStr)
	if err1 != nil || err2 != nil {
		m.logger.Error(fmt.Sprintf("[esclient][GetResponse]err during query marshal=%s %s", err1, err2))
	}
	m.logger.Info(fmt.Sprintf("[esclient]Final ESQuery=%s", string(queryJs)))

	searchService := m.client.Search(elastic2.GamesIndex).SearchSource(searchSource)
	searchResult, err := searchService.Do(ctx)
	if err != nil {
		m.logger.Error("cannot perform search")
		return nil, err
	}

	result := make([]*protos.Game, 0, len(searchResult.Hits.Hits))
	for _, g := range searchResult.Hits.Hits {
		var game *protos.Game
		if err := json.Unmarshal(g.Source, &game); err != nil {
			m.logger.Error("error while unmarshalling search result hit")
			continue
		}
		result = append(result, game)
	}

	return result, nil
}

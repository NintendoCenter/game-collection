package service

import (
	"context"
	"fmt"

	"NintendoCenter/game-collection/internal/manager"
	"NintendoCenter/game-collection/internal/protos"
	"go.uber.org/zap"
)

type SearchService interface {
	SearchGames(ctx context.Context, filter *protos.FindGameRequest) ([]*protos.Game, error)
}

type GameService struct {
	l         *zap.Logger
	manager   *manager.GameManager
	esManager *manager.ElasticManager
}

func NewGameService(logger *zap.Logger, manager *manager.GameManager, esManager *manager.ElasticManager) *GameService {
	return &GameService{l: logger, manager: manager, esManager: esManager}
}

func (m *GameService) SaveGame(ctx context.Context, game *protos.Game) error {
	var err error
	if existed, _ := m.manager.Find(game.Id); existed != nil {
		m.l.Info(fmt.Sprintf("game '%s' updated", game.Title))
		err = m.manager.UpdateGame(game.Id, game)
	} else {
		m.l.Info(fmt.Sprintf("game '%s' saved", game.Title))
		err = m.manager.SaveGame(game)
	}

	if err != nil {
		return err
	}

	return m.esManager.IndexGame(ctx, game)
}

func (m *GameService) SearchGames(ctx context.Context, filter *protos.FindGameRequest) ([]*protos.Game, error) {
	return m.manager.SearchGames(ctx, filter)
}

func (m *GameService) GetGame(_ context.Context, id string) (*protos.Game, error) {
	return m.manager.Find(id)
}

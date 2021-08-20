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

func NewGameService(logger *zap.Logger, manager *manager.GameManager) *GameService {
	return &GameService{l: logger, manager: manager}
}

func (m *GameService) SaveGame(ctx context.Context, game *protos.Game) error {
	var err error
	if existed, _ := m.manager.Find(ctx, game.Id); existed != nil {
		m.l.Info(fmt.Sprintf("game '%s' updated", game.Title))
		err = m.manager.UpdateGame(ctx, game.Id, game)
	} else {
		m.l.Info(fmt.Sprintf("game '%s' saved", game.Title))
		err = m.manager.SaveGame(ctx, game)
	}

	if err != nil {
		return err
	}

	if m.esManager != nil {
		return m.esManager.IndexGame(ctx, game)
	}

	return nil
}

func (m *GameService) SearchGames(ctx context.Context, filter *protos.FindGameRequest) ([]*protos.Game, error) {
	return m.manager.SearchGames(ctx, filter)
}

func (m *GameService) GetGame(ctx context.Context, id string) (*protos.Game, error) {
	return m.manager.Find(ctx, id)
}

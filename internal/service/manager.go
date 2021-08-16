package service

import (
	"fmt"

	"NintendoCenter/game-collection/internal/manager"
	"NintendoCenter/game-collection/internal/protos"
	"go.uber.org/zap"
)

type GameService struct {
	l       *zap.Logger
	manager *manager.GameManager
}

func NewGameService(logger *zap.Logger, manager *manager.GameManager) *GameService {
	return &GameService{l: logger, manager: manager}
}

func (m *GameService) SaveGame(game *protos.Game) error {
	if existed, _ := m.manager.Find(game.Id); existed != nil {
		if existed == game {
			m.l.Info("game already exists in a same state. Skipped.")
			return nil
		}
		m.l.Info(fmt.Sprintf("game '%s' updated", game.Title))
		return m.manager.UpdateGame(game.Id, game)
	}
	m.l.Info(fmt.Sprintf("game '%s' saved", game.Title))
	return m.manager.SaveGame(game)
}

func (m *GameService) SearchGames(name string) ([]*protos.Game, error) {
	return m.manager.SearchByName(name)
}

func (m *GameService) GetGame(id string) (*protos.Game, error) {
	return m.manager.Find(id)
}
package service

import (
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
		return m.manager.UpdateGame(game.Id, game)
	}
	return m.manager.SaveGame(game)
}

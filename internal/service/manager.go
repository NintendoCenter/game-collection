package service

import (
	"fmt"

	"NintendoCenter/game-collection/internal/protos"
	"go.uber.org/zap"
)

type GameManager struct {
	l *zap.Logger
	// TODO: mgo
}

func NewGameManager(logger *zap.Logger) *GameManager {
	return &GameManager{l: logger}
}

func (m *GameManager) SaveGame(game *protos.Game) error {
	m.l.Info(fmt.Sprintf("Got game '%s' for saving", game.Title))
	return nil
}

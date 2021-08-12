package main

import (
	"NintendoCenter/game-collection/config"
	"NintendoCenter/game-collection/internal/infrastructure"
	"NintendoCenter/game-collection/internal/providers/grpc_server"
	"NintendoCenter/game-collection/internal/providers/logger"
	"NintendoCenter/game-collection/internal/service"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func buildContainer() (*dig.Container, error) {
	container := dig.New()

	constructors := []interface{}{
		func (cfg *config.Config) (*zap.Logger, error) {
			return logger.New(cfg.LogLevel)
		},
		func (i *infrastructure.CollectionServer, cfg *config.Config) (*grpc_server.GrpcServer, error) {
			return grpc_server.New(i, cfg)
		},
		func (l *zap.Logger) *service.GameManager {
			return service.NewGameManager(l)
		},
		func (m *service.GameManager) *infrastructure.CollectionServer {
			return infrastructure.NewCollectionServer(m)
		},
		config.NewConfig,
	}

	for _, c := range constructors {
		if err := container.Provide(c); err != nil {
			return nil, err
		}
	}

	return container, nil
}

package providers

import (
	"NintendoCenter/game-collection/config"
	"NintendoCenter/game-collection/internal/infrastructure"
	"NintendoCenter/game-collection/internal/manager"
	"NintendoCenter/game-collection/internal/providers/grpc_server"
	"NintendoCenter/game-collection/internal/providers/logger"
	"NintendoCenter/game-collection/internal/providers/mongo"
	"NintendoCenter/game-collection/internal/queue/consumer"
	"NintendoCenter/game-collection/internal/service"
	"github.com/globalsign/mgo"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func BuildContainer() (*dig.Container, error) {
	container := dig.New()

	constructors := []interface{}{
		func (cfg *config.Config) (*zap.Logger, error) {
			return logger.New(cfg.LogLevel)
		},
		func (i *infrastructure.CollectionServer, cfg *config.Config) (*grpc_server.GrpcServer, error) {
			return grpc_server.New(i, cfg)
		},
		func (l *zap.Logger, m *manager.GameManager) *service.GameService {
			return service.NewGameService(l, m)
		},
		func (m *service.GameService) *infrastructure.CollectionServer {
			return infrastructure.NewCollectionServer(m)
		},
		func (cfg *config.Config) (*mgo.Database, error) {
			return mongo.New(cfg.MongoConnection)
		},
		func (db *mgo.Database) (*manager.GameManager, error) {
			return manager.NewGameManager(db)
		},
		func (s *service.GameService, cfg *config.Config, l *zap.Logger) (*consumer.GameConsumer, error) {
			return consumer.NewGameConsumer(cfg.GamesTopic, cfg.QueueAddr, s, l)
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

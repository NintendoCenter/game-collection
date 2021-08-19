package consumer

import (
	"context"
	"encoding/json"

	"NintendoCenter/game-collection/internal/protos"
	"NintendoCenter/game-collection/internal/service"
	"github.com/nsqio/go-nsq"
	"go.uber.org/zap"
)

type GameConsumer struct {
	service *service.GameService
	consumer *nsq.Consumer
	logger *zap.Logger
	address string
	ctx *context.Context
}

const DefaultChannel = "default"

func NewGameConsumer(topic string, queueAddress string, service *service.GameService, logger *zap.Logger) (*GameConsumer, error) {
	consumer, err := nsq.NewConsumer(topic, DefaultChannel, nsq.NewConfig())
	if err != nil {
		return nil, err
	}

	gameConsumer := &GameConsumer{
		consumer: consumer,
		logger:   logger,
		address: queueAddress,
		service: service,
	}

	consumer.AddHandler(gameConsumer)

	return gameConsumer, nil
}

func (c *GameConsumer) Start(ctx context.Context) error {
	c.ctx = &ctx
	return c.consumer.ConnectToNSQD(c.address)
}

func (c *GameConsumer) Stop() {
	c.logger.Info("stopping consumer")
	c.consumer.Stop()
}

func (c *GameConsumer) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	var game protos.Game
	if err := json.Unmarshal(m.Body, &game); err != nil {
		return err
	}

	return c.service.SaveGame(*c.ctx, &game)
}

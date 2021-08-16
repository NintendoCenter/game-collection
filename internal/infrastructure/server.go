package infrastructure

import (
	"context"

	"NintendoCenter/game-collection/internal/protos"
	"NintendoCenter/game-collection/internal/service"
)

type CollectionServer struct {
	m *service.GameService
}

func NewCollectionServer(m *service.GameService) *CollectionServer {
	return &CollectionServer{
		m: m,
	}
}

func (c *CollectionServer) SaveGame(ctx context.Context, game *protos.Game) (*protos.Game, error) {
	return game, c.m.SaveGame(game)
}

func (c *CollectionServer) SaveOffer(ctx context.Context, request *protos.SaveOfferRequest) (*protos.Offer, error) {
	return nil, nil
}

func (c *CollectionServer) GetGame(ctx context.Context, request *protos.GetGameRequest) (*protos.Game, error) {
	return nil, nil
}

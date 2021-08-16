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

func (c *CollectionServer) GetGame(ctx context.Context, r *protos.GetGameRequest) (*protos.Game, error) {
	return c.m.GetGame(r.GetId())
}

func (c *CollectionServer) FindGame(ctx context.Context, r *protos.FindGameRequest) (*protos.FindGameResponse, error) {
	games, err := c.m.SearchGames(r.GetTitle())
	if err != nil {
		return nil, err
	}

	return &protos.FindGameResponse{Games: games}, nil
}

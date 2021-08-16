package grpc_server

import (
	"fmt"
	"net"

	"NintendoCenter/game-collection/config"
	"NintendoCenter/game-collection/internal/infrastructure"
	"NintendoCenter/game-collection/internal/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	gs *grpc.Server
	i *infrastructure.CollectionServer
	cfg *config.Config
}

func New(i *infrastructure.CollectionServer, cfg *config.Config) (*GrpcServer, error) {
	gs := grpc.NewServer()
	protos.RegisterGameCollectionServer(gs, i)
	if cfg.EnableReflection {
		reflection.Register(gs)
	}

	return &GrpcServer{
		gs: gs,
		i:  i,
		cfg: cfg,
	}, nil
}

func (s *GrpcServer) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.GrpcPort))
	if err != nil {
		return err
	}
	return s.gs.Serve(l)
}

func (s *GrpcServer)Stop() {
	s.gs.Stop()
}

package main

import (
	"fmt"

	"NintendoCenter/game-collection/config"
	"NintendoCenter/game-collection/internal/infrastructure"
	"NintendoCenter/game-collection/internal/providers"
	"NintendoCenter/game-collection/internal/providers/grpc_server"
	"go.uber.org/zap"
)

func main() {
	t1, _ := zap.NewProduction()
	container, err := providers.BuildContainer()
	if err != nil {
		t1.Fatal("cannot build dependencies", zap.Error(err))
	}

	err = container.Invoke(func(cfg *config.Config, server *grpc_server.GrpcServer, logger *zap.Logger, cs *infrastructure.CollectionServer) {
		err := server.Run()
		if err != nil {
			logger.Fatal(fmt.Sprintf("cannot start listener on port %d", cfg.GrpcPort), zap.Error(err))
		}
	})

	if err != nil {
		t1.Fatal("container cannot invoke it's dependencies", zap.Error(err))
	}
}

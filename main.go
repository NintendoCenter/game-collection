package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"NintendoCenter/game-collection/internal/providers"
	"NintendoCenter/game-collection/internal/providers/grpc_server"
	"NintendoCenter/game-collection/internal/queue/consumer"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {
	t1, _ := zap.NewProduction()
	container, err := providers.BuildContainer()
	if err != nil {
		t1.Fatal("cannot build dependencies", zap.Error(err))
	}

	gr, ctx := errgroup.WithContext(context.Background())
	errStopped := errors.New("service stopped")

	err = container.Invoke(func(s *grpc_server.GrpcServer, c *consumer.GameConsumer, logger *zap.Logger) {
		defer c.Stop()
		defer s.Stop()
		gr.Go(func() error {
			if err := c.Start(ctx); err != nil {
				logger.Fatal("error on starting consumer")
				return err
			}
			return s.Run()
		})

		gr.Go(func() error {
			signals := make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
			defer signal.Stop(signals)
			defer s.Stop()

			select {
			case <- ctx.Done():
				return ctx.Err()
			case <- signals:
				return errStopped
			}
		})

		if err := gr.Wait(); err != nil && err != errStopped {
			t1.Fatal("caught an error. Terminating", zap.Error(err))
		}

		t1.Info("signal to quit")
	})

	if err != nil {
		t1.Fatal("container cannot invoke it's dependencies", zap.Error(err))
	}
}

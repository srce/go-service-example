package api

import (
	"context"
	"fmt"

	"github.com/dzyanis/go-service-example/pkg/logger"
)

type Boot struct {
	log *logger.Logger
	api *Server
}

func NewBoot(log *logger.Logger, api *Server) Boot {
	return Boot{
		log: log,
		api: api,
	}
}

func (b Boot) Name() string {
	return "http server"
}

func (b Boot) Init(ctx context.Context) error {
	if _, err := RegisterHandlers(b.log, b.api); err != nil {
		return fmt.Errorf("registering handlers: %w", err)
	}

	go func() {
		b.log.Printf("listening on %s\n", b.api.Addr)
		if err := b.api.ListenAndServe(); err != nil {
			b.log.Printf("%s: %v", b.Name(), err)
		}
	}()

	return nil
}

func (b Boot) Stop() error {
	if err := b.api.Close(); err != nil {
		return fmt.Errorf("stopping: %w", err)
	}

	return nil
}

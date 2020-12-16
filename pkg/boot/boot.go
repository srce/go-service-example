package boot

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dzyanis/go-service-example/pkg/logger"
)

var (
	ErrServiceFailed = errors.New("service failed")
	ErrPanic         = errors.New("boot panic")
)

type booter interface {
	Name() string
	Init(ctx context.Context) error
	Stop() error
}

type Runner struct {
	ctx      context.Context
	log      *logger.Logger
	services []booter
}

func NewRunner(ctx context.Context, log *logger.Logger) *Runner {
	return &Runner{
		ctx:      ctx,
		log:      log,
		services: make([]booter, 0),
	}
}

func (r *Runner) safeInit(b booter) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w: %v", ErrPanic, r)
		}
	}()

	return b.Init(r.ctx)
}

func (r *Runner) Try(b booter, attempts int, delay time.Duration) error {
	for ; attempts > 0; attempts-- {
		if err := r.safeInit(b); err != nil {
			r.log.Printf("service %s: %v", b.Name(), err)
			time.Sleep(delay)

			continue
		}

		r.log.Printf("service %s is working", b.Name())
		r.services = append(r.services, b)
		return nil
	}

	return fmt.Errorf("%s: %w", b.Name(), ErrServiceFailed)
}

func (r *Runner) Once(b booter) error {
	if err := b.Init(r.ctx); err != nil {
		return fmt.Errorf("service %s failed: %w", b.Name(), err)
	}
	r.services = append(r.services, b)
	return nil
}

func (r *Runner) WaitForTermination() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	r.log.Print("call shut down")

	for _, service := range r.services {
		if err := service.Stop(); err != nil {
			r.log.Printf("%s stopping failed: %v", service.Name(), err)
		} else {
			r.log.Printf("%s was stopped", service.Name())
		}
	}
}

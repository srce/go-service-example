package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dzyanis/go-service-example/internal/transactions"
	"github.com/dzyanis/go-service-example/internal/users"
	"github.com/dzyanis/go-service-example/internal/wallets"
)

type Server struct {
	*http.Server
	users        *users.Controller
	wallets      *wallets.Controller
	transactions *transactions.Controller
}

type Config struct {
	Host    string        `env:"API_HOST"`
	Port    int           `env:"API_PORT"`
	Timeout time.Duration `env:"API_TIMEOUT"`
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func NewServer(cfg Config, users *users.Controller,
	wallets *wallets.Controller, transactions *transactions.Controller) *Server {
	return &Server{
		Server: &http.Server{
			Addr:         cfg.Addr(),
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
		},
		users:        users,
		wallets:      wallets,
		transactions: transactions,
	}
}

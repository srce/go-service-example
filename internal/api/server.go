package api

import (
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	*http.Server
}

type Config struct {
	Host    string        `env:"API_HOST"`
	Port    int           `env:"API_PORT"`
	Timeout time.Duration `env:"API_TIMEOUT"`
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func NewServer(cfg Config) *Server {
	return &Server{
		Server: &http.Server{
			Addr:         cfg.Addr(),
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
		},
	}
}

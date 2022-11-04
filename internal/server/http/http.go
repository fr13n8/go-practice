package http

import (
	"fmt"

	"github.com/fr13n8/go-practice/internal/config"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	http *fiber.App
	host string
	port string
}

func (s *Server) Run(cfg *config.HttpConfig, initHandlers func(app *fiber.App)) error {
	srv := fiber.New(
		fiber.Config{
			AppName:      cfg.ServerName,
			WriteTimeout: cfg.WriteTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			BodyLimit:    cfg.MaxRequestBodySize,
		},
	)

	s.host = cfg.Host
	s.port = cfg.Port
	s.http = srv
	initHandlers(srv)
	return s.http.Listen(fmt.Sprintf("%s:%s", s.host, s.port))
}

func (s *Server) Shutdown() error {
	return s.http.Shutdown()
}

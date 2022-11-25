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

func NewServer(cfg *config.HttpConfig) *Server {
	return &Server{
		http: fiber.New(
			fiber.Config{
				AppName:      cfg.ServerName,
				WriteTimeout: cfg.WriteTimeout,
				ReadTimeout:  cfg.ReadTimeout,
				BodyLimit:    cfg.MaxRequestBodySize,
			}),
		host: cfg.Host,
		port: cfg.Port,
	}
}

func (s *Server) Run(initHandlers func(app *fiber.App)) error {
	initHandlers(s.http)
	return s.http.Listen(fmt.Sprintf("%s:%s", s.host, s.port))
}

func (s *Server) Shutdown() error {
	return s.http.Shutdown()
}

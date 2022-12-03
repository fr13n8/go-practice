package http

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
				WriteTimeout: cfg.WriteTimeout * time.Second,
				ReadTimeout:  cfg.ReadTimeout * time.Second,
				BodyLimit:    cfg.MaxRequestBodySize,
			}),
		host: cfg.Host,
		port: cfg.Port,
	}
}

func (s *Server) Run(initHandlers func(app *fiber.App)) <-chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	initHandlers(s.http)

	go func() {
		if err := s.http.Listen(fmt.Sprintf("%s:%s", s.host, s.port)); err != nil {
			log.Fatalf("Error while running server: %s", err.Error())
		}
	}()

	return quit
}

func (s *Server) ShutdownGracefully() {
	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		// Release resources like Database connections
		cancel()
	}()

	shutdownChan := make(chan error, 1)
	go func() { shutdownChan <- s.http.Shutdown() }()

	select {
	case <-timeout.Done():
		log.Fatal("Server Shutdown Timed out before shutdown.")
	case err := <-shutdownChan:
		if err != nil {
			log.Fatal("Error while shutting down server", err)
		} else {
			fmt.Println("Server Shutdown Successful")
		}
	}
}

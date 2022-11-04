package app

import (
	"fmt"
	"github.com/fr13n8/go-practice/pkg/redis"
	"os"
	"os/signal"
	"syscall"

	"github.com/fr13n8/go-practice/internal/config"
	httpHandler "github.com/fr13n8/go-practice/internal/delivary/http"
	"github.com/fr13n8/go-practice/internal/repository"
	httpServer "github.com/fr13n8/go-practice/internal/server/http"
	"github.com/fr13n8/go-practice/internal/services"
	"github.com/fr13n8/go-practice/pkg/database"
)

func Run() {
	cfg := config.NewConfig()
	fmt.Println("Config:", cfg)
	db := database.NewDb(cfg.Database)
	cache, err := redis.NewRedis()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	repo := repository.NewRepository(db, cache)
	appServices := services.NewServices(repo)
	handlers := httpHandler.NewHandler(appServices)
	srv := new(httpServer.Server)

	go func() {
		if err := srv.Run(&cfg.HTTP, handlers.Init); err != nil {
			fmt.Printf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(); err != nil {
		fmt.Printf("error occurred while shutting down http server: %s\n", err.Error())
	}

	fmt.Println("Server stopped")
}

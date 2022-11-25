package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	httpServer "github.com/fr13n8/go-practice/internal/delivary/http"
	httpHandler "github.com/fr13n8/go-practice/internal/delivary/http/handler"

	grpcServer "github.com/fr13n8/go-practice/internal/delivary/grpc"
	grpcHandler "github.com/fr13n8/go-practice/internal/delivary/grpc/handler"

	"github.com/fr13n8/go-practice/pkg/redis"

	"github.com/fr13n8/go-practice/internal/config"
	"github.com/fr13n8/go-practice/internal/repository"
	"github.com/fr13n8/go-practice/internal/services"
	"github.com/fr13n8/go-practice/pkg/database"
)

func RunHttp() {
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
	srv := httpServer.NewServer(&cfg.HTTP)

	go func() {
		if err := srv.Run(handlers.Init); err != nil {
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

func RunGrpc() {
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
	handlers := grpcHandler.NewHandler(appServices)
	srv := grpcServer.NewServer(&cfg.GRPC)
	go func() {
		if err := srv.Run(handlers.Init); err != nil {
			fmt.Printf("error occurred while running grpc server: %s\n", err.Error())
		}
	}()
	fmt.Printf("gRPC server is running on port %s", cfg.GRPC.Port)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	srv.Stop()

	fmt.Println("Server stopped")
}

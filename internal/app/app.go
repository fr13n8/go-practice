package app

import (
	"fmt"
	"os"

	httpServer "github.com/fr13n8/go-practice/internal/delivary/http"
	httpHandler "github.com/fr13n8/go-practice/internal/delivary/http/handler"
	"github.com/fr13n8/go-practice/internal/delivary/http/middlewares"

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

	middlewares.RegisterPromethesMetrics()
	quit := srv.Run(handlers.Init)
	<-quit

	fmt.Println("Shutting down http server...")
	srv.ShutdownGracefully()
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

	quit := srv.Run(handlers.Init)
	<-quit

	fmt.Println("Shutting down gRPC server...")
	srv.ShutdownGracefully()
}

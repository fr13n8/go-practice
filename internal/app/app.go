package app

import (
	"fmt"
	"os"

	httpServer "github.com/fr13n8/go-practice/internal/delivary/http"
	httpHandler "github.com/fr13n8/go-practice/internal/delivary/http/handler"
	"github.com/opentracing/opentracing-go"

	grpcServer "github.com/fr13n8/go-practice/internal/delivary/grpc"
	grpcHandler "github.com/fr13n8/go-practice/internal/delivary/grpc/handler"

	jaegerTracer "github.com/fr13n8/go-practice/pkg/jaeger"

	"github.com/fr13n8/go-practice/pkg/redis"
	"github.com/fr13n8/go-practice/pkg/utils"

	"github.com/fr13n8/go-practice/internal/config"
	"github.com/fr13n8/go-practice/internal/repository"
	"github.com/fr13n8/go-practice/internal/services"
	"github.com/fr13n8/go-practice/pkg/database"
)

func RunHttp() {
	configPath := utils.GetConfigPath(os.Getenv("config"))
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	db := database.NewDb(&cfg.Database)
	cache, err := redis.NewRedis(&cfg.Redis)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	repo := repository.NewRepository(db, cache)
	appServices := services.NewServices(repo)
	handlers := httpHandler.NewHandler(appServices)
	srv := httpServer.NewServer(&cfg.HTTP)

	tracer, closer, err := jaegerTracer.InitJaeger(&cfg.Jaeger, "http")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Jaeger tracer started")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	fmt.Println("Opentracing connected")

	quit := srv.Run(handlers.Init)
	<-quit

	fmt.Println("Shutting down http server...")
	srv.ShutdownGracefully()
}

func RunGrpc() {
	configPath := utils.GetConfigPath(os.Getenv("config"))
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	db := database.NewDb(&cfg.Database)
	cache, err := redis.NewRedis(&cfg.Redis)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	repo := repository.NewRepository(db, cache)
	appServices := services.NewServices(repo)
	handlers := grpcHandler.NewHandler(appServices)
	srv := grpcServer.NewServer(&cfg.GRPC)

	tracer, closer, err := jaegerTracer.InitJaeger(&cfg.Jaeger, "grpc")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Jaeger tracer started")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	fmt.Println("Opentracing connected")

	quit := srv.Run(handlers.Init)
	<-quit

	fmt.Println("Shutting down gRPC server...")
	srv.ShutdownGracefully()
}

package grpc

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fr13n8/go-practice/internal/config"
	pb "github.com/fr13n8/go-practice/pkg/grpc/v1/gen"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Deps struct {
	TaskHandler pb.TaskServiceServer
}

type Server struct {
	gRpcServer *grpc.Server
	host       string
	port       string
}

func NewServer(cfg *config.GrpcConfig) *Server {
	return &Server{
		gRpcServer: grpc.NewServer(
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle: cfg.MaxConnectionIdle * time.Minute,
				Timeout:           cfg.Timeout * time.Second,
				MaxConnectionAge:  cfg.MaxConnectionAge * time.Minute,
				Time:              cfg.Timeout * time.Minute,
			}),
			// grpc.UnaryInterceptor(im.GrpcLogger),
			grpc.ChainUnaryInterceptor(
				grpc_ctxtags.UnaryServerInterceptor(),
				grpc_prometheus.UnaryServerInterceptor,
				grpcrecovery.UnaryServerInterceptor(),
			),
		),
		host: cfg.Host,
		port: cfg.Port,
	}
}

func (s *Server) Run(initHandlers func(app *grpc.Server)) <-chan os.Signal {
	addr := s.host + ":" + s.port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	initHandlers(s.gRpcServer)

	go func() {
		if err := s.gRpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return quit
}

func (s *Server) ShutdownGracefully() {
	s.gRpcServer.GracefulStop()
	fmt.Println("gRPC Server Shutdown Successful")
}

package grpc

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/fr13n8/go-practice/pkg/grpc/v1"

	"github.com/fr13n8/go-practice/internal/config"
	"google.golang.org/grpc"
)

type Deps struct {
	TaskHandler pb.TaskServiceServer
}

type Server struct {
	srv  *grpc.Server
	host string
	port string
}

func NewServer(cfg *config.GrpcConfig) *Server {
	return &Server{
		srv:  grpc.NewServer(),
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

	initHandlers(s.srv)

	go func() {
		if err := s.srv.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return quit
}

func (s *Server) ShutdownGracefully() {
	s.srv.GracefulStop()
	fmt.Println("gRPC Server Shutdown Successful")
}

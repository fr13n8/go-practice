package grpc

import (
	"net"

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

func (s *Server) Run(initHandlers func(app *grpc.Server)) error {
	addr := s.host + ":" + s.port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	initHandlers(s.srv)

	return s.srv.Serve(listener)
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}

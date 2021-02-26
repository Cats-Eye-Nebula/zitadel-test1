package auth

import (
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/internal/auth/repository"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/pkg/grpc/auth"
	"google.golang.org/grpc"
)

var _ auth.AuthServiceServer = (*Server)(nil)

const (
	authName = "Auth-API"
)

type Server struct {
	auth.UnimplementedAuthServiceServer
	command *command.Commands
	query   *query.Queries
	repo    repository.Repository
}

type Config struct {
	Repository eventsourcing.Config
}

func CreateServer(command *command.Commands, query *query.Queries, authRepo repository.Repository) *Server {
	return &Server{
		command: command,
		query:   query,
		repo:    authRepo,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	auth.RegisterAuthServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return authName
}

func (s *Server) MethodPrefix() string {
	return auth.AuthService_MethodPrefix
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return auth.AuthService_AuthMethods
}

func (s *Server) RegisterGateway() server.GatewayFunc {
	return auth.RegisterAuthServiceHandlerFromEndpoint
}

func (s *Server) GatewayPathPrefix() string {
	return "/auth/v1"
}

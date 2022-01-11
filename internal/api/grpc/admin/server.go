package admin

import (
	"google.golang.org/grpc"

	"github.com/caos/zitadel/internal/admin/repository"
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/server"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/pkg/grpc/admin"
)

const (
	adminName = "Admin-API"
)

var _ admin.AdminServiceServer = (*Server)(nil)

type Server struct {
	admin.UnimplementedAdminServiceServer
	command         *command.Commands
	query           *query.Queries
	iam             repository.IAMRepository
	administrator   repository.AdministratorRepository
	repo            repository.Repository
	users           repository.UserRepository
	iamDomain       string
	assetsAPIDomain string
}

type Config struct {
	Repository eventsourcing.Config
}

func CreateServer(command *command.Commands, query *query.Queries, repo repository.Repository, iamDomain, assetsAPIDomain string) *Server {
	return &Server{
		command:         command,
		query:           query,
		iam:             repo,
		administrator:   repo,
		repo:            repo,
		users:           repo,
		iamDomain:       iamDomain,
		assetsAPIDomain: assetsAPIDomain,
	}
}

func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	admin.RegisterAdminServiceServer(grpcServer, s)
}

func (s *Server) AppName() string {
	return adminName
}

func (s *Server) MethodPrefix() string {
	return admin.AdminService_MethodPrefix
}

func (s *Server) AuthMethods() authz.MethodMapping {
	return admin.AdminService_AuthMethods
}

func (s *Server) RegisterGateway() server.GatewayFunc {
	return admin.RegisterAdminServiceHandlerFromEndpoint
}

func (s *Server) GatewayPathPrefix() string {
	return "/admin/v1"
}
